package response

type successResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func Success(message string, data any) successResponse {
	return successResponse{Message: message, Data: data}
}

func Error(message string) errorResponse {
	return errorResponse{Message: message}
}
