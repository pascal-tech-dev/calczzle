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
		// --- success: guide / core cases ---
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
			name:  "simple exponentiation",
			input: "2 ^ 3",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "^"},
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
			name:  "left associative subtraction",
			input: "3 - 4 - 5",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "-"},
				{Type: TokenNumber, Value: "5"},
				{Type: TokenOperator, Value: "-"},
			},
		},
		{
			name:  "left associative multiplication and division",
			input: "3 * 4 / 2",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "*"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "/"},
			},
		},
		{
			name:  "nested parentheses",
			input: "3 * (4 + 5)",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenNumber, Value: "5"},
				{Type: TokenOperator, Value: "+"},
				{Type: TokenOperator, Value: "*"},
			},
		},
		{
			name:  "double parentheses around number",
			input: "((3))",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
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
			name:  "standalone percentage",
			input: "100%",
			want: []Token{
				{Type: TokenNumber, Value: "100"},
				{Type: TokenPercentage, Value: "%"},
			},
		},
		{
			name:  "decimal numbers",
			input: "3.5 + 2.25",
			want: []Token{
				{Type: TokenNumber, Value: "3.5"},
				{Type: TokenNumber, Value: "2.25"},
				{Type: TokenOperator, Value: "+"},
			},
		},

		// --- success: unary operators ---
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
			name:  "unary minus lower than exponentiation",
			input: "-2 ^ 2",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "^"},
				{Type: TokenOperator, Value: "u-"},
			},
		},
		{
			name:  "unary minus higher than multiplication",
			input: "-2 * 3",
			want: []Token{
				{Type: TokenNumber, Value: "2"},
				{Type: TokenOperator, Value: "u-"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "*"},
			},
		},
		{
			name:  "unary minus inside parentheses",
			input: "(-3)",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "u-"},
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
			name:  "double plus treats second as unary",
			input: "3++4",
			want: []Token{
				{Type: TokenNumber, Value: "3"},
				{Type: TokenNumber, Value: "4"},
				{Type: TokenOperator, Value: "+"},
			},
		},

		// --- errors ---
		{
			name:    "consecutive operators",
			input:   "3 + * 4",
			wantErr: true,
		},
		{
			name:    "leading binary operator",
			input:   "* 3",
			wantErr: true,
		},
		{
			name:    "missing operator between numbers",
			input:   "3 4",
			wantErr: true,
		},
		{
			name:    "number after closing parenthesis",
			input:   "(3 + 4) 5",
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
			name:    "empty parentheses",
			input:   "()",
			wantErr: true,
		},
		{
			name:    "empty sqrt arguments",
			input:   "sqrt()",
			wantErr: true,
		},
		{
			name:    "bare function without call",
			input:   "sqrt",
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
		{
			name:    "lone plus",
			input:   "+",
			wantErr: true,
		},
		{
			name:    "lone minus",
			input:   "-",
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

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		_, err := ToPostfix(nil)
		if err == nil {
			t.Fatal("ToPostfix(nil) error = nil, want error")
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		_, err := ToPostfix([]Token{})
		if err == nil {
			t.Fatal("ToPostfix([]) error = nil, want error")
		}
	})
}

// TestToPostfixFromTokens exercises the parser without the tokenizer,
// so invalid token sequences are attributed only to ToPostfix.
func TestToPostfixFromTokens(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tokens  []Token
		want    []Token
		wantErr bool
	}{
		{
			name: "simple subtraction tokens",
			tokens: []Token{
				{Type: TokenNumber, Value: "10", Position: 0},
				{Type: TokenOperator, Value: "-", Position: 2},
				{Type: TokenNumber, Value: "3", Position: 4},
			},
			want: []Token{
				{Type: TokenNumber, Value: "10"},
				{Type: TokenNumber, Value: "3"},
				{Type: TokenOperator, Value: "-"},
			},
		},
		{
			name: "function then left paren then number then right paren",
			tokens: []Token{
				{Type: TokenFunction, Value: "sqrt", Position: 0},
				{Type: TokenLeftParenthesis, Value: "(", Position: 4},
				{Type: TokenNumber, Value: "9", Position: 5},
				{Type: TokenRightParenthesis, Value: ")", Position: 6},
			},
			want: []Token{
				{Type: TokenNumber, Value: "9"},
				{Type: TokenFunction, Value: "sqrt"},
			},
		},
		{
			name: "unknown operator symbol",
			tokens: []Token{
				{Type: TokenNumber, Value: "1", Position: 0},
				{Type: TokenOperator, Value: "@", Position: 1},
				{Type: TokenNumber, Value: "2", Position: 2},
			},
			wantErr: true,
		},
		{
			name: "unexpected left parenthesis after number",
			tokens: []Token{
				{Type: TokenNumber, Value: "1", Position: 0},
				{Type: TokenLeftParenthesis, Value: "(", Position: 1},
			},
			wantErr: true,
		},
		{
			name: "function left on stack without call",
			tokens: []Token{
				{Type: TokenFunction, Value: "sqrt", Position: 0},
				{Type: TokenNumber, Value: "4", Position: 5},
			},
			wantErr: true,
		},
		{
			name: "invalid token type",
			tokens: []Token{
				{Type: TokenType(99), Value: "?", Position: 0},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ToPostfix(tt.tokens)
			if tt.wantErr {
				if err == nil {
					t.Fatal("ToPostfix(...) error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ToPostfix(...) unexpected error: %v", err)
			}
			assertTokens(t, got, tt.want)
		})
	}
}
