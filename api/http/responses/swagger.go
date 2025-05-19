package responses

type AuthSuccessResponse struct {
	// in: body
	Code    int64            `json:"code" example:"1000"`
	Message string           `json:"message" example:"token successfully generated"`
	Data    AuthResponseData `json:"data"`
}

// ErrorResponse is the standard envelope for error responses.
// swagger:response ErrorResponse
type ErrorResponse struct {
	// in: body

	// Code is your internal API status code, e.g. 1002
	Code int64 `json:"code" example:"0000"`
	// Message is a human‑readable description, e.g. "invalid json request body"
	Message string `json:"message"`
	Data    string `json:"data"`
}
