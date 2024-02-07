package core

type SuccessResponse struct {
	Result interface{} `json:"result"`
}

type ErrorResponse struct {
	Error interface{} `json:"error"`
}
