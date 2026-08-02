package controller

import "strconv"

// EvaluateResponse is the JSON body returned on successful evaluation.
type EvaluateResponse struct {
	Result float64 `json:"result"`
}

// errorBody is the public error envelope returned by calculator endpoints.
type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// cleanFloat removes IEEE-754 binary noise (e.g. 0.10500000000000001 → 0.105)
// while keeping enough significant digits for calculator results.
func cleanFloat(v float64) float64 {
	formatted := strconv.FormatFloat(v, 'g', 15, 64)
	cleaned, err := strconv.ParseFloat(formatted, 64)
	if err != nil {
		return v
	}
	return cleaned
}
