package expression

import (
	"errors"

	"pascal-tech-dev/calczzle-backend/pkg"
)

// ToPostfix converts infix tokens to postfix using the shunting-yard algorithm.
func ToPostfix(tokens []Token) ([]Token, error) {
	if len(tokens) == 0 {
		return nil, errors.New("empty expression")
	}

	output := pkg.NewQueue[Token]()
	stack := pkg.NewStack[Token]()
	expectOperand := true

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			if !expectOperand {
				return nil, errors.New("unexpected number")
			}
			output.Enqueue(token)
			expectOperand = false

		case TokenFunction:
			if !expectOperand {
				return nil, errors.New("unexpected function")
			}
			if !IsSupportedFunction(token.Value) {
				return nil, ErrUnsupportedFunction
			}
			stack.Push(token)

		case TokenLeftParenthesis:
			if !expectOperand {
				return nil, errors.New("unexpected left parenthesis")
			}
			stack.Push(token)

		case TokenRightParenthesis:
			if expectOperand {
				return nil, errors.New("unexpected right parenthesis")
			}
			if err := popUntilLeftParen(output, stack); err != nil {
				return nil, err
			}
			expectOperand = false

		case TokenPercentage:
			if expectOperand {
				return nil, errors.New("unexpected percentage")
			}
			output.Enqueue(token)

		case TokenOperator:
			if expectOperand {
				if err := handleUnaryOperator(token, stack); err != nil {
					return nil, err
				}
				continue
			}
			if err := handleBinaryOperator(token, output, stack); err != nil {
				return nil, err
			}
			expectOperand = true

		default:
			return nil, errors.New("invalid token")
		}
	}

	if expectOperand {
		return nil, errors.New("expression ends unexpectedly")
	}

	for !stack.IsEmpty() {
		top, _ := stack.Pop()
		switch top.Type {
		case TokenLeftParenthesis:
			return nil, errors.New("mismatched parentheses")
		case TokenFunction:
			return nil, errors.New("mismatched parentheses")
		case TokenOperator:
			output.Enqueue(top)
		default:
			return nil, errors.New("invalid token on operator stack")
		}
	}

	return drainQueue(output), nil
}

func handleUnaryOperator(token Token, stack *pkg.Stack[Token]) error {
	switch token.Value {
	case "+":
		// Unary plus is a no-op.
		return nil
	case "-":
		stack.Push(Token{
			Type:     TokenOperator,
			Value:    UnaryMinusSymbol,
			Position: token.Position,
		})
		return nil
	default:
		return errors.New("unexpected operator")
	}
}

func handleBinaryOperator(token Token, output *pkg.Queue[Token], stack *pkg.Stack[Token]) error {
	op, ok := LookupOperator(token.Value)
	if !ok {
		return errors.New("unknown operator")
	}

	for {
		top, ok := stack.Peek()
		if !ok || !shouldPopOperator(top, op) {
			break
		}
		stack.Pop()
		output.Enqueue(top)
	}

	stack.Push(token)
	return nil
}

func shouldPopOperator(top Token, incoming Operator) bool {
	topOp, ok := operatorFromStackToken(top)
	if !ok {
		return false
	}

	if topOp.Associativity == LeftAssociative {
		return topOp.Precedence >= incoming.Precedence
	}
	return topOp.Precedence > incoming.Precedence
}

func operatorFromStackToken(token Token) (Operator, bool) {
	if token.Type != TokenOperator {
		return Operator{}, false
	}
	return LookupOperator(token.Value)
}

func popUntilLeftParen(output *pkg.Queue[Token], stack *pkg.Stack[Token]) error {
	found := false

	for !stack.IsEmpty() {
		top, _ := stack.Pop()

		if top.Type == TokenLeftParenthesis {
			found = true
			break
		}
		if top.Type == TokenFunction {
			return errors.New("mismatched parentheses")
		}
		output.Enqueue(top)
	}

	if !found {
		return errors.New("mismatched parentheses")
	}

	if top, ok := stack.Peek(); ok && top.Type == TokenFunction {
		fn, _ := stack.Pop()
		output.Enqueue(fn)
	}

	return nil
}

func drainQueue(q *pkg.Queue[Token]) []Token {
	result := make([]Token, 0, q.Len())
	for !q.IsEmpty() {
		item, _ := q.Dequeue()
		result = append(result, item)
	}
	return result
}
