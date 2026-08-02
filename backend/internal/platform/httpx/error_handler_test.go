package httpx_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"

	"pascal-tech-dev/calczzle-backend/internal/calculator/expression"
	"pascal-tech-dev/calczzle-backend/internal/platform/httpx"
)

func TestErrorHandlerMapping(t *testing.T) {
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
			name:       "invalid expression sentinel",
			err:        expression.ErrInvalidExpression,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_EXPRESSION",
		},
		{
			name:       "plain syntax error",
			err:        errors.New("mismatched parentheses"),
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_EXPRESSION",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := fiber.New(fiber.Config{ErrorHandler: httpx.ErrorHandler})
			app.Get("/boom", func(ctx fiber.Ctx) error {
				return tt.err
			})

			req := httptest.NewRequest(http.MethodGet, "/boom", nil)
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

			var got httpx.ErrorResponse
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.Error.Code != tt.wantCode {
				t.Fatalf("code = %q, want %q", got.Error.Code, tt.wantCode)
			}
			if got.Error.Message == "" {
				t.Fatal("expected non-empty message")
			}
		})
	}
}
