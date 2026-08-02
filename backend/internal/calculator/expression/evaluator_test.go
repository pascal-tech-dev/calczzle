package expression

import (
	"errors"
	"math"
	"testing"
)

func TestEvaluatePostfixPipeline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr error
	}{
		{name: "simple addition", input: "3 + 4", want: 7},
		{name: "mixed operators from guide", input: "3 + 4 - 7 * 2", want: -7},
		{name: "parentheses change order", input: "(3 + 4 - 7) * 2", want: 0},
		{name: "simple exponentiation", input: "2 ^ 3", want: 8},
		{name: "right associative exponentiation", input: "2 ^ 3 ^ 2", want: 512},
		{name: "sqrt function", input: "sqrt(81)", want: 9},
		{name: "percentage", input: "100 * 15%", want: 15},
		{name: "unary minus", input: "-3 + 4", want: 1},
		{name: "unary minus lower than exponentiation", input: "-2 ^ 2", want: -4},
		{name: "unary minus after operator", input: "2 ^ -3", want: 0.125},
		{name: "unary minus higher than multiplication", input: "-2 * 3", want: -6},
		{name: "nested sqrt", input: "sqrt((9 + 7) * 4)", want: 8},
		{name: "division by zero", input: "10 / 0", wantErr: ErrDivisionByZero},
		{name: "negative square root", input: "sqrt(-4)", wantErr: ErrInvalidSquareRoot},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokens, err := Tokenize(tt.input)
			if err != nil {
				t.Fatalf("Tokenize(%q): %v", tt.input, err)
			}
			postfix, err := ToPostfix(tokens)
			if err != nil {
				t.Fatalf("ToPostfix(%q): %v", tt.input, err)
			}

			got, err := EvaluatePostfix(postfix)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("EvaluatePostfix(%q) error = %v, want %v", tt.input, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("EvaluatePostfix(%q) unexpected error: %v", tt.input, err)
			}
			assertFloat(t, got, tt.want)
		})
	}
}

func TestEvaluatePostfix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tokens  []Token
		want    float64
		wantErr error
	}{
		{
			name: "guide postfix mixed operators",
			tokens: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenOperator, Value: "-"},
			},
			want: -7,
		},
		{
			name: "percentage unary",
			tokens: []Token{
				{Type: TokenNumber, Value: "15"},
				{Type: TokenPercentage, Value: "%"},
			},
			want: 0.15,
		},
		{
			name: "unary minus token",
			tokens: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: UnaryMinusSymbol},
			},
			want: -3,
		},
		{
			name:    "empty tokens",
			tokens:  nil,
			wantErr: ErrInvalidExpression,
		},
		{
			name: "stack underflow on binary",
			tokens: []Token{
				{Type: TokenNumber, Value: "1"},
				{Type: TokenOperator, Value: "+"},
			},
			wantErr: ErrInvalidExpression,
		},
		{
			name: "extra operands",
			tokens: []Token{
				{Type: TokenNumber, Value: "1"},
				{Type: TokenNumber, Value: "2"},
			},
			wantErr: ErrInvalidExpression,
		},
		{
			name: "division by zero",
			tokens: []Token{
				{Type: TokenNumber, Value: "10"},
				{Type: TokenNumber, Value: "0"},
				{Type: TokenOperator, Value: "/"},
			},
			wantErr: ErrDivisionByZero,
		},
		{
			name: "negative square root",
			tokens: []Token{
				{Type: TokenNumber, Value: "-4"},
				{Type: TokenFunction, Value: "sqrt"},
			},
			wantErr: ErrInvalidSquareRoot,
		},
		{
			name: "unknown operator",
			tokens: []Token{
				{Type: TokenNumber, Value: "1"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "@"},
			},
			wantErr: ErrInvalidExpression,
		},
		{
			name: "parenthesis in postfix is invalid",
			tokens: []Token{
				{Type: TokenLeftParenthesis, Value: "("},
			},
			wantErr: ErrInvalidExpression,
		},
		{
			name: "non-finite power",
			tokens: []Token{
				{Type: TokenNumber, Value: "-1"},
				{Type: TokenNumber, Value: "0.5"},
				{Type: TokenOperator, Value: "^"},
			},
			wantErr: ErrInvalidExpression,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := EvaluatePostfix(tt.tokens)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("EvaluatePostfix(...) error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("EvaluatePostfix(...) unexpected error: %v", err)
			}
			assertFloat(t, got, tt.want)
		})
	}
}

func assertFloat(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("result = %v, want %v", got, want)
	}
}
