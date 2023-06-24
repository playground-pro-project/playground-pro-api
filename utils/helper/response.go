package helper

type RequestResponse struct {
	Code       int         `json:"code,omitempty"`
	Message    string      `json:"message"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func ResponseFormat(code int, message string, data interface{}, pagination interface{}) RequestResponse {
	result := RequestResponse{
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	return result
}

func ErrorResponse(err string) RequestResponse {
	return RequestResponse{
		Error: err,
	}
}

func SuccessResponse(data interface{}, msg string) RequestResponse {
	return RequestResponse{
		Data:    data,
		Message: msg,
	}
}
