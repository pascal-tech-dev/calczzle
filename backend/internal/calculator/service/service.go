package service

// Service orchestrates expression evaluation.
type Service struct{}

// New constructs a calculator Service.
func New() *Service {
	return &Service{}
}

// Evaluate evaluates an arithmetic expression.
// Temporary stub until the expression engine is connected.
func (s *Service) Evaluate(_ string) (float64, error) {
	return 42, nil
}
