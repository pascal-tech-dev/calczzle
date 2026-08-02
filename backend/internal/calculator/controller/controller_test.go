package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"pascal-tech-dev/calczzle-backend/internal/calculator/controller"
	"pascal-tech-dev/calczzle-backend/internal/calculator/expression"
	"pascal-tech-dev/calczzle-backend/internal/platform/httpx"
)

type fakeEvaluator struct {
	result float64
	err    error
	called string
}

func (f *fakeEvaluator) Evaluate(expression string) (float64, error) {
	f.called = expression
	return f.result, f.err
}

func TestEvaluateSuccess(t *testing.T) {
	t.Parallel()

	fake := &fakeEvaluator{result: 42}
	app := fiber.New()
	controller.New(fake).Register(app)

	body := []byte(`{"expression":"3 + 4"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var got controller.EvaluateResponse
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Result != 42 {
		t.Fatalf("result = %v, want 42", got.Result)
	}
	if fake.called != "3 + 4" {
		t.Fatalf("called with %q, want %q", fake.called, "3 + 4")
	}
}

func TestEvaluateEmptyExpression(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	controller.New(&fakeEvaluator{}).Register(app)

	body := []byte(`{"expression":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var got map[string]map[string]string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got["error"]["code"] != "EMPTY_EXPRESSION" {
		t.Fatalf("code = %q, want EMPTY_EXPRESSION", got["error"]["code"])
	}
}

func TestEvaluateInvalidJSON(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	controller.New(&fakeEvaluator{}).Register(app)

	body := []byte(`{`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	var got map[string]map[string]string
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got["error"]["code"] != "INVALID_REQUEST" {
		t.Fatalf("code = %q, want INVALID_REQUEST", got["error"]["code"])
	}
}

func TestEvaluateDomainErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "division by zero",
			err:        expression.ErrDivisionByZero,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "DIVISION_BY_ZERO",
		},
		{
			name:       "invalid square root",
			err:        expression.ErrInvalidSquareRoot,
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   "INVALID_SQUARE_ROOT",
		},
		{
			name:       "unsupported function",
			err:        expression.ErrUnsupportedFunction,
			wantStatus: http.StatusBadRequest,
			wantCode:   "UNSUPPORTED_FUNCTION",
		},
		{
			name:       "invalid expression",
			err:        errors.New("mismatched parentheses"),
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_EXPRESSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New(fiber.Config{ErrorHandler: httpx.ErrorHandler})
			controller.New(&fakeEvaluator{err: tt.err}).Register(app)

			body := []byte(`{"expression":"1 + 1"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			raw, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}

			var got map[string]map[string]string
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got["error"]["code"] != tt.wantCode {
				t.Fatalf("code = %q, want %q", got["error"]["code"], tt.wantCode)
			}
		})
	}
}
