package controller

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
