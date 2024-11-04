package responses

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewBaseResponse(status int, message string) *BaseResponse {
	return &BaseResponse{
		Status:  status,
		Message: message,
	}
}
