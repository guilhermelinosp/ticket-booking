package responses

import (
	"ticket-booking/entities"
)

// EventResponse is the base response for all responses.
type EventResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    []*entities.Event `json:"data,omitempty"`
}

// NewEventResponse creates a new instance of BaseResponse.
func NewEventResponse(status int, message string, data []*entities.Event) *EventResponse {
	return &EventResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
