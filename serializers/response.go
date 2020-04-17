package serializers

// Response format...
type Response struct {
	Code    int         `json:"code"`
	Body    interface{} `json:"body"`
	Message string      `json:"message"`
}

// SerializeError ..
func SerializeError(code int, message string) interface{} {
	return NewError(code, message)
}

// SerializeResponse ..
func SerializeResponse(code int, body interface{}, message string) interface{} {
	return NewResponse(code, body, message)
}

// NewError error message format
func NewError(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// NewResponse response foramt
func NewResponse(code int, body interface{}, message string) Response {
	return Response{
		Code:    code,
		Body:    body,
		Message: message,
	}
}
