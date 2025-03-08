package api

type Valid interface {
	Validate() error
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func GetErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}

func Ok() StatusResponse {
	return StatusResponse{
		Status: "ok",
	}
}

type MessageResponse struct {
	Message string `json:"message"`
}

func GetMessageResponse(err error) MessageResponse {
	return MessageResponse{
		Message: err.Error(),
	}
}
