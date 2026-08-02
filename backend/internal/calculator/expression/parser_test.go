package expression

import "testing"

func TestToPostfix(t *testing.T) {
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
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
			},
		},
		{
			name:  "mixed operators from guide",
			input: "3 + 4 - 7 * 2",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenOperator, Value: "-"},
			},
		},
		{
			name:  "parentheses change order",
			input: "(3 + 4 - 7) * 2",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "*"},
			},
		},
		{
			name:  "right associative exponentiation",
			input: "2 ^ 3 ^ 2",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenOperator, Value: "^"},
			},
		},
		{
			name:  "sqrt function",
			input: "sqrt(81)",
			want: []Token{
				{Type: TokenNumber, Value: "81"},
				{Type: TokenFunction, Value: "sqrt"},
			},
		},
		{
			name:  "percentage",
			input: "100 * 15%",
			want: []Token{
				{Type: TokenNumber, Value: "100"},
				{Type: TokenNumber, Value: "15"},
				{Type: TokenPercentage, Value: "%"},
				{Type: TokenOperator, Value: "*"},
			},
		},
		{
			name:  "unary minus",
			input: "-3 + 4",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "u-"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
			},
		},
		{
			name:  "unary minus after operator",
			input: "2 ^ -3",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "u-"},
				{Type: TokenOperator, Value: "^"},
			},
		},
		{
			name:  "unary plus is ignored",
			input: "+3 + 4",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
			},
		},
		{
			name:  "nested sqrt expression",
			input: "sqrt((9 + 7) * 4)",
			want: []Token{
				{Type: TokenNumber, Value: "9"},
				{Type: TokenNumber, Value: "7"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenFunction, Value: "sqrt"},
			},
		},
		{
			name:    "consecutive operators",
			input:   "3 + * 4",
			wantErr: true,
		},
		{
			name:    "unmatched left parenthesis",
			input:   "(3 + 4",
			wantErr: true,
		},
		{
			name:    "unmatched right parenthesis",
			input:   "3 + 4)",
			wantErr: true,
		},
		{
			name:    "unsupported function",
			input:   "foo(1)",
			wantErr: true,
		},
		{
			name:    "percentage at start",
			input:   "%15",
			wantErr: true,
		},
		{
			name:    "percentage after operator",
			input:   "3 + %",
			wantErr: true,
		},
		{
			name:    "trailing operator",
			input:   "3 +",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokens, err := Tokenize(tt.input)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Fatalf("Tokenize(%q): %v", tt.input, err)
			}

			got, err := ToPostfix(tokens)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ToPostfix(%q) error = nil, want error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToPostfix(%q) unexpected error: %v", tt.input, err)
			}
			assertTokens(t, got, tt.want)
		})
	}
}

func TestToPostfixEmptyTokens(t *testing.T) {
	t.Parallel()

	_, err := ToPostfix(nil)
	if err == nil {
		t.Fatal("ToPostfix(nil) error = nil, want error")
	}
}
