package expression

import "errors"

var (
	// ErrDivisionByZero is returned when dividing by zero.
	ErrDivisionByZero = errors.New("division by zero")

	// ErrInvalidSquareRoot is returned when taking the square root of a negative number.
	ErrInvalidSquareRoot = errors.New("square root of negative number")

	// ErrUnsupportedFunction is returned when an expression uses an unknown function.
	ErrUnsupportedFunction = errors.New("unsupported function")

	// ErrInvalidExpression is returned for malformed postfix evaluation
	// (stack underflow, extra operands, unknown tokens, non-finite results).
	ErrInvalidExpression = errors.New("invalid expression")
)
