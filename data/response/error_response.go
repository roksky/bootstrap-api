package response

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message,omitempty"`
}
