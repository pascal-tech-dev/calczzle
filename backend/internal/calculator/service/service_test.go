package service_test

import (
	"errors"
	"math"
	"testing"

	"pascal-tech-dev/calczzle-backend/internal/calculator/expression"
	"pascal-tech-dev/calczzle-backend/internal/calculator/service"
)

func TestEvaluate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr error
		anyErr  bool
	}{
		{name: "simple addition", input: "3 + 4", want: 7},
		{name: "mixed operators from guide", input: "3 + 4 - 7 * 2", want: -7},
		{name: "parentheses change order", input: "(3 + 4 - 7) * 2", want: 0},
		{name: "simple exponentiation", input: "2 ^ 3", want: 8},
		{name: "right associative exponentiation", input: "2 ^ 3 ^ 2", want: 512},
		{name: "sqrt function", input: "sqrt(81)", want: 9},
		{name: "percentage", input: "100 * 15%", want: 15},
		{name: "unary minus", input: "-3 + 4", want: 1},
		{name: "division by zero", input: "10 / 0", wantErr: expression.ErrDivisionByZero},
		{name: "negative square root", input: "sqrt(-4)", wantErr: expression.ErrInvalidSquareRoot},
		{name: "unsupported function", input: "foo(1)", wantErr: expression.ErrUnsupportedFunction},
		{name: "consecutive operators", input: "3 + * 4", anyErr: true},
		{name: "unmatched parenthesis", input: "(3 + 4", anyErr: true},
		{name: "empty expression", input: "", anyErr: true},
	}

	svc := service.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := svc.Evaluate(tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Evaluate(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}
			if tt.anyErr {
				if err == nil {
					t.Fatalf("Evaluate(%q) error = nil, want error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Evaluate(%q) unexpected error: %v", tt.input, err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("Evaluate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEvaluateImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ service.Evaluator = service.New()
}
