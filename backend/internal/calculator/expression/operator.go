package expression

type Associativity int

const (
	LeftAssociative Associativity = iota
	RightAssociative
)

type Operator struct {
	Symbol        string
	Precedence    int
	Associativity Associativity
}

var operators = map[string]Operator{
	"+": {Symbol: "+", Precedence: 1, Associativity: LeftAssociative},
	"-": {Symbol: "-", Precedence: 1, Associativity: LeftAssociative},
	"*": {Symbol: "*", Precedence: 2, Associativity: LeftAssociative},
	"/": {Symbol: "/", Precedence: 2, Associativity: LeftAssociative},
	"^": {Symbol: "^", Precedence: 3, Associativity: RightAssociative},
}

// LookupOperator returns the operator definition for a symbol.
func LookupOperator(symbol string) (Operator, bool) {
	op, ok := operators[symbol]
	return op, ok
}
