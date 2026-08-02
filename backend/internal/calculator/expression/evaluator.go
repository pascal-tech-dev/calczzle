package expression

import (
	"math"
	"strconv"

	"pascal-tech-dev/calczzle-backend/pkg"
)

// EvaluatePostfix evaluates a postfix token sequence and returns the result.
func EvaluatePostfix(tokens []Token) (float64, error) {
	if len(tokens) == 0 {
		return 0, ErrInvalidExpression
	}

	stack := pkg.NewStack[float64]()

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, ErrInvalidExpression
			}
			stack.Push(value)

		case TokenOperator:
			if token.Value == UnaryMinusSymbol {
				if err := applyUnary(stack, func(x float64) (float64, error) {
					return -x, nil
				}); err != nil {
					return 0, err
				}
				continue
			}
			if err := applyBinary(stack, token.Value); err != nil {
				return 0, err
			}

		case TokenPercentage:
			if err := applyUnary(stack, func(x float64) (float64, error) {
				return x / 100, nil
			}); err != nil {
				return 0, err
			}

		case TokenFunction:
			if token.Value != "sqrt" {
				return 0, ErrInvalidExpression
			}
			if err := applyUnary(stack, func(x float64) (float64, error) {
				if x < 0 {
					return 0, ErrInvalidSquareRoot
				}
				return math.Sqrt(x), nil
			}); err != nil {
				return 0, err
			}

		default:
			return 0, ErrInvalidExpression
		}
	}

	if stack.Len() != 1 {
		return 0, ErrInvalidExpression
	}

	result, _ := stack.Pop()
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrInvalidExpression
	}
	return result, nil
}

func applyBinary(stack *pkg.Stack[float64], op string) error {
	right, left, err := popTwo(stack)
	if err != nil {
		return err
	}

	var result float64
	switch op {
	case "+":
		result = left + right
	case "-":
		result = left - right
	case "*":
		result = left * right
	case "/":
		if right == 0 {
			return ErrDivisionByZero
		}
		result = left / right
	case "^":
		result = math.Pow(left, right)
		if math.IsNaN(result) || math.IsInf(result, 0) {
			return ErrInvalidExpression
		}
	default:
		return ErrInvalidExpression
	}

	stack.Push(result)
	return nil
}

func applyUnary(stack *pkg.Stack[float64], fn func(float64) (float64, error)) error {
	x, err := popOne(stack)
	if err != nil {
		return err
	}
	result, err := fn(x)
	if err != nil {
		return err
	}
	stack.Push(result)
	return nil
}

func popOne(stack *pkg.Stack[float64]) (float64, error) {
	value, ok := stack.Pop()
	if !ok {
		return 0, ErrInvalidExpression
	}
	return value, nil
}

func popTwo(stack *pkg.Stack[float64]) (right, left float64, err error) {
	right, err = popOne(stack)
	if err != nil {
		return 0, 0, err
	}
	left, err = popOne(stack)
	if err != nil {
		return 0, 0, err
	}
	return right, left, nil
}
