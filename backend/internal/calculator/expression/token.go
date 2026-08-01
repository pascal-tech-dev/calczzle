package expression

type TokenType int

const (
	TokenNumber TokenType = iota
	TokenOperator
	TokenLeftParenthesis
	TokenRightParenthesis
	TokenFunction
	TokenPercentage
)

type Token struct {
	Type     TokenType
	Value    string
	Position int
}
