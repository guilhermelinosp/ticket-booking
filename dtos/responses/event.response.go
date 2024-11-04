package responses

import (
	"ticket-booking/entities"
)

type EventResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    []*entities.Event `json:"data,omitempty"`
}

func NewEventResponse(status int, message string, data []*entities.Event) *EventResponse {
	return &EventResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}


