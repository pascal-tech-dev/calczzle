package expression

import "testing"

func TestLookupOperator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		symbol        string
		wantOK        bool
		wantPrecedence int
		wantAssoc     Associativity
	}{
		{
			name:           "addition",
			symbol:         "+",
			wantOK:         true,
			wantPrecedence: 1,
			wantAssoc:      LeftAssociative,
		},
		{
			name:           "subtraction",
			symbol:         "-",
			wantOK:         true,
			wantPrecedence: 1,
			wantAssoc:      LeftAssociative,
		},
		{
			name:           "multiplication",
			symbol:         "*",
			wantOK:         true,
			wantPrecedence: 2,
			wantAssoc:      LeftAssociative,
		},
		{
			name:           "division",
			symbol:         "/",
			wantOK:         true,
			wantPrecedence: 2,
			wantAssoc:      LeftAssociative,
		},
		{
			name:           "exponentiation",
			symbol:         "^",
			wantOK:         true,
			wantPrecedence: 3,
			wantAssoc:      RightAssociative,
		},
		{
			name:   "percentage is not a binary operator",
			symbol: "%",
			wantOK:  false,
		},
		{
			name:   "unknown symbol",
			symbol: "@",
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := LookupOperator(tt.symbol)
			if ok != tt.wantOK {
				t.Fatalf("LookupOperator(%q) ok = %v, want %v", tt.symbol, ok, tt.wantOK)
			}
			if !tt.wantOK {
				return
			}
			if got.Symbol != tt.symbol {
				t.Fatalf("Symbol = %q, want %q", got.Symbol, tt.symbol)
			}
			if got.Precedence != tt.wantPrecedence {
				t.Fatalf("Precedence = %d, want %d", got.Precedence, tt.wantPrecedence)
			}
			if got.Associativity != tt.wantAssoc {
				t.Fatalf("Associativity = %v, want %v", got.Associativity, tt.wantAssoc)
			}
		})
	}
}
