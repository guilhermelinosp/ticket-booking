package responses

import "ticket-booking/entities"

// TicketResponse is the base response for all responses.
type TicketResponse struct {
	Status  int                `json:"status"`
	Message string             `json:"message"`
	Data    []*entities.Ticket `json:"data,omitempty"`
}

// NewTicketResponse creates a new instance of BaseResponse.
func NewTicketResponse(status int, message string, data []*entities.Ticket) *TicketResponse {
	return &TicketResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
