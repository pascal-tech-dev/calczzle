package service

import "pascal-tech-dev/calczzle-backend/internal/calculator/expression"

// Service orchestrates expression evaluation.
type Service struct{}

// New constructs a calculator Service.
func New() *Service {
	return &Service{}
}

// Evaluate evaluates an arithmetic expression.
func (s *Service) Evaluate(input string) (float64, error) {
	tokens, err := expression.Tokenize(input)
	if err != nil {
		return 0, err
	}

	postfix, err := expression.ToPostfix(tokens)
	if err != nil {
		return 0, err
	}

	return expression.EvaluatePostfix(postfix)
}
