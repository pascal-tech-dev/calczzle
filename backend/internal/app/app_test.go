package app_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pascal-tech-dev/calczzle-backend/internal/app"
	"pascal-tech-dev/calczzle-backend/internal/config"
)

func TestBodyLimitRejectsOversizedRequest(t *testing.T) {
	t.Parallel()

	server := app.New(config.Config{Port: "8080"})

	// Build a JSON body larger than MaxRequestBodyBytes.
	padding := strings.Repeat("1", app.MaxRequestBodyBytes)
	body := []byte(`{"expression":"` + padding + `"}`)
	if len(body) <= app.MaxRequestBodyBytes {
		t.Fatalf("test body length %d must exceed limit %d", len(body), app.MaxRequestBodyBytes)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// fiber.App.Test rejects oversized bodies before producing a response.
	_, err := server.Test(req)
	if err == nil {
		t.Fatal("expected error for oversized body, got nil")
	}
	if !strings.Contains(err.Error(), "body size exceeds the given limit") {
		t.Fatalf("error = %v, want body size limit error", err)
	}
}

func TestBodyLimitAllowsNormalRequest(t *testing.T) {
	t.Parallel()

	server := app.New(config.Config{Port: "8080"})

	body := []byte(`{"expression":"3 + 4"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := server.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, raw)
	}
}
