package controller

// EvaluateRequest is the JSON body for POST /api/v1/evaluate.
type EvaluateRequest struct {
	Expression string `json:"expression"`
}
