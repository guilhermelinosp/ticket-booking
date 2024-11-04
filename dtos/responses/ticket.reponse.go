package responses

import "ticket-booking/entities"

type TicketResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    []*entities.Ticket `json:"data,omitempty"`
}

func NewTicketResponse(status int, message string, data []*entities.Ticket) *TicketResponse {
	return &TicketResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
