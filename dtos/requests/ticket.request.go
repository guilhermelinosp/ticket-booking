package requests

import "github.com/go-playground/validator/v10"

type TicketRequest struct {
	TicketId string `json:"ticket_id" validate:"required"`
	EventId  string `json:"event_id" validate:"required"`
	UserId   string `json:"user_id" validate:"required"`
}

func NewTicketRequest(ticketId, eventId, userId string) *TicketRequest {
	return &TicketRequest{
		TicketId: ticketId,
		EventId:  eventId,
		UserId:   userId,
	}
}

func (t *TicketRequest) Validate() error {
	return validator.New().Struct(t)
}
