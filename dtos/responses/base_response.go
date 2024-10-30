package responses

type BaseResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*EventResponse `json:"data,omitempty"`
}

func NewBaseResponse(status, message string, data []*EventResponse) *BaseResponse {
	return &BaseResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
