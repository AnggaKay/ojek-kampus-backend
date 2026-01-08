package dto

type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"` // Validation errors
}

func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(code, message string) Response {
	return Response{
		Success: false,
		Message: message,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}

func ValidationErrorResponse(fields map[string]string) Response {
	return Response{
		Success: false,
		Message: "Validation failed",
		Error: &ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: "Invalid input data",
			Fields:  fields,
		},
	}
}
