package service_test

import (
	"testing"

	"pascal-tech-dev/calczzle-backend/internal/calculator/service"
)

func TestEvaluateStub(t *testing.T) {
	t.Parallel()

	svc := service.New()
	got, err := svc.Evaluate("3 + 4")
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if got != 42 {
		t.Fatalf("result = %v, want 42", got)
	}
}
