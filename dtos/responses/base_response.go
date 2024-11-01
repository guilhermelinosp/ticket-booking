package responses

import (
	"ticket-booking/entities"
)

// BaseResponse is the base response for all responses.
type BaseResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    []*entities.Event `json:"data,omitempty"`
}

// NewBaseResponse creates a new instance of BaseResponse.
func NewBaseResponse(status int, message string, data []*entities.Event) *BaseResponse {
	return &BaseResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
