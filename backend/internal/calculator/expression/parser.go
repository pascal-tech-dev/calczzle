package expression

import "errors"

const unaryMinusSymbol = "u-"

var unaryMinusOp = Operator{
	Symbol:        unaryMinusSymbol,
	Precedence:    4,
	Associativity: RightAssociative,
}

// ToPostfix converts infix tokens to postfix using the shunting-yard algorithm.
func ToPostfix(tokens []Token) ([]Token, error) {
	if len(tokens) == 0 {
		return nil, errors.New("empty expression")
	}

	output := make([]Token, 0, len(tokens))
	stack := make([]Token, 0)
	expectOperand := true

	for _, token := range tokens {
		switch token.Type {
		case TokenNumber:
			if !expectOperand {
				return nil, errors.New("unexpected number")
			}
			output = append(output, token)
			expectOperand = false

		case TokenFunction:
			if !expectOperand {
				return nil, errors.New("unexpected function")
			}
			if !IsSupportedFunction(token.Value) {
				return nil, errors.New("unsupported function")
			}
			stack = append(stack, token)

		case TokenLeftParenthesis:
			if !expectOperand {
				return nil, errors.New("unexpected left parenthesis")
			}
			stack = append(stack, token)

		case TokenRightParenthesis:
			if expectOperand {
				return nil, errors.New("unexpected right parenthesis")
			}
			if err := popUntilLeftParen(&output, &stack); err != nil {
				return nil, err
			}
			expectOperand = false

		case TokenPercentage:
			if expectOperand {
				return nil, errors.New("unexpected percentage")
			}
			output = append(output, token)

		case TokenOperator:
			if expectOperand {
				if err := handleUnaryOperator(token, &stack); err != nil {
					return nil, err
				}
				continue
			}
			if err := handleBinaryOperator(token, &output, &stack); err != nil {
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

	for i := len(stack) - 1; i >= 0; i-- {
		top := stack[i]
		switch top.Type {
		case TokenLeftParenthesis:
			return nil, errors.New("mismatched parentheses")
		case TokenFunction:
			return nil, errors.New("mismatched parentheses")
		case TokenOperator:
			output = append(output, top)
		default:
			return nil, errors.New("invalid token on operator stack")
		}
	}

	return output, nil
}

func handleUnaryOperator(token Token, stack *[]Token) error {
	switch token.Value {
	case "+":
		// Unary plus is a no-op.
		return nil
	case "-":
		*stack = append(*stack, Token{
			Type:     TokenOperator,
			Value:    unaryMinusSymbol,
			Position: token.Position,
		})
		return nil
	default:
		return errors.New("unexpected operator")
	}
}

func handleBinaryOperator(token Token, output *[]Token, stack *[]Token) error {
	op, ok := LookupOperator(token.Value)
	if !ok {
		return errors.New("unknown operator")
	}

	for len(*stack) > 0 {
		top := (*stack)[len(*stack)-1]
		if !shouldPopOperator(top, op) {
			break
		}
		*output = append(*output, top)
		*stack = (*stack)[:len(*stack)-1]
	}

	*stack = append(*stack, token)
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
	if token.Value == unaryMinusSymbol {
		return unaryMinusOp, true
	}
	return LookupOperator(token.Value)
}

func popUntilLeftParen(output *[]Token, stack *[]Token) error {
	found := false

	for len(*stack) > 0 {
		top := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]

		if top.Type == TokenLeftParenthesis {
			found = true
			break
		}
		if top.Type == TokenFunction {
			return errors.New("mismatched parentheses")
		}
		*output = append(*output, top)
	}

	if !found {
		return errors.New("mismatched parentheses")
	}

	if len(*stack) > 0 && (*stack)[len(*stack)-1].Type == TokenFunction {
		fn := (*stack)[len(*stack)-1]
		*stack = (*stack)[:len(*stack)-1]
		*output = append(*output, fn)
	}

	return nil
}
