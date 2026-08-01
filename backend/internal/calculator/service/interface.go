package service

// Evaluator evaluates arithmetic expressions.
type Evaluator interface {
	Evaluate(expression string) (float64, error)
}
