package expression

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{
			name:  "simple addition",
			input: "3 + 4",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
			},
		},
		{
			name:  "mixed operators",
			input: "3 + 4 - 7 * 2",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "2"},
			},
		},
		{
			name:  "decimal number",
			input: "3.14",
			want: []Token{
				{Type: TokenNumber, Value: "3.14"},
			},
		},
		{
			name:  "parentheses",
			input: "(3 + 4) * 2",
			want: []Token{
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "2"},
			},
		},
		{
			name:  "exponentiation",
			input: "2 ^ 3",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenNumber, Value: "3"},
			},
		},
		{
			name:  "sqrt function",
			input: "sqrt(81)",
			want: []Token{
				{Type: TokenFunction, Value: "sqrt"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "81"},
				{Type: TokenRightParenthesis, Value: ")"},
			},
		},
		{
			name:  "function name is lowercased",
			input: "SQRT(9)",
			want: []Token{
				{Type: TokenFunction, Value: "sqrt"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "9"},
				{Type: TokenRightParenthesis, Value: ")"},
			},
		},
		{
			name:  "percentage",
			input: "100 * 15%",
			want: []Token{
				{Type: TokenNumber, Value: "100"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "15"},
				{Type: TokenOperator, Value: "%"},
			},
		},
		{
			name:  "unicode operators are normalized",
			input: "6 × 4 ÷ 2 − 1",
			want: []Token{
				{Type: TokenNumber, Value: "6"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "/"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "1"},
			},
		},
		{
			name:  "whitespace is ignored",
			input: "  1   +   2  ",
			want: []Token{
				{Type: TokenNumber, Value: "1"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "2"},
			},
		},
		{
			name:  "nested parentheses and mixed precedence",
			input: "((3 + 4) - 7) * 2 / (1 + 1)",
			want: []Token{
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "/"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "1"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "1"},
				{Type: TokenRightParenthesis, Value: ")"},
			},
		},
		{
			name:  "no spaces between tokens",
			input: "12.5*(3+4.25)/2^3",
			want: []Token{
				{Type: TokenNumber, Value: "12.5"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4.25"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "/"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenNumber, Value: "3"},
			},
		},
		{
			name:  "right-associative power chain",
			input: "2 ^ 3 ^ 2",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenNumber, Value: "2"},
			},
		},
		{
			name:  "unary minus and plus are just operators",
			input: "-3 + +4 - -5",
			want: []Token{
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "5"},
			},
		},
		{
			name:  "nested sqrt with expression argument",
			input: "sqrt((9 + 7) * 4) + sqrt(1)",
			want: []Token{
				{Type: TokenFunction, Value: "sqrt"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "9"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenFunction, Value: "sqrt"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "1"},
				{Type: TokenRightParenthesis, Value: ")"},
			},
		},
		{
			name:  "percentage of decimal expression",
			input: "(12.5 + 7.5) * 20%",
			want: []Token{
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "12.5"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "7.5"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "20"},
				{Type: TokenOperator, Value: "%"},
			},
		},
		{
			name:  "leading decimal is accepted when digits follow",
			input: ".5 + 1.",
			want: []Token{
				{Type: TokenNumber, Value: ".5"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "1."},
			},
		},
		{
			name:  "tabs and newlines count as whitespace",
			input: "1\t+\n2\r\n*\t3",
			want: []Token{
				{Type: TokenNumber, Value: "1"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "3"},
			},
		},
		{
			name:  "identifier with underscore and digits",
			input: "sqrt2_fn(16)",
			want: []Token{
				{Type: TokenFunction, Value: "sqrt2_fn"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "16"},
				{Type: TokenRightParenthesis, Value: ")"},
			},
		},
		{
			name:  "mixed unicode and ascii operators",
			input: "10 × (2 + 3) ÷ 5 − 1 ^ 2",
			want: []Token{
				{Type: TokenNumber, Value: "10"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenRightParenthesis, Value: ")"},
				{Type: TokenOperator, Value: "/"},
				{Type: TokenNumber, Value: "5"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "1"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenNumber, Value: "2"},
			},
		},
		{
			name:  "unmatched parentheses still tokenize",
			input: "(3 + 4",
			want: []Token{
				{Type: TokenLeftParenthesis, Value: "("},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
			},
		},
		{
			name:  "consecutive operators still tokenize",
			input: "3 + * 4",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "4"},
			},
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "bare decimal point",
			input:   ".",
			wantErr: true,
		},
		{
			name:    "multiple decimal points",
			input:   "1.2.3",
			wantErr: true,
		},
		{
			name:    "multiple decimal points like date",
			input:   "3.12.23",
			wantErr: true,
		},
		{
			name:    "three consecutive decimal points",
			input:   "3...",
			wantErr: true,
		},
		{
			name:    "double decimal in longer number",
			input:   "12.34.56",
			wantErr: true,
		},
		{
			name:    "leading and middle decimal points",
			input:   ".3.14",
			wantErr: true,
		},
		{
			name:    "trailing double decimal",
			input:   "5..",
			wantErr: true,
		},
		{
			name:    "multiple decimals after operator",
			input:   "10 + 3.12.23",
			wantErr: true,
		},
		{
			name:    "multiple decimals inside parentheses",
			input:   "(3.12.23)",
			wantErr: true,
		},
		{
			name:    "multiple decimals in function argument",
			input:   "sqrt(9.0.0)",
			wantErr: true,
		},
		{
			name:    "invalid character",
			input:   "3 @ 4",
			wantErr: true,
		},
		{
			name:    "invalid character after valid prefix",
			input:   "sqrt(9) + $1",
			wantErr: true,
		},
		{
			name:    "trailing incomplete decimal in larger expression",
			input:   "2 + .",
			wantErr: true,
		},
		{
			name:    "number with two dots mid expression",
			input:   "(1.2.3 + 4)",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Tokenize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Tokenize(%q) error = nil, want error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("Tokenize(%q) unexpected error: %v", tt.input, err)
			}
			assertTokens(t, got, tt.want)
		})
	}
}

func assertTokens(t *testing.T, got, want []Token) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("token count = %d, want %d\ngot:  %#v\nwant: %#v", len(got), len(want), got, want)
	}

	for i := range want {
		if got[i].Type != want[i].Type || got[i].Value != want[i].Value {
			t.Fatalf("token[%d] = {%v %q}, want {%v %q}",
				i, got[i].Type, got[i].Value, want[i].Type, want[i].Value)
		}
	}
}
