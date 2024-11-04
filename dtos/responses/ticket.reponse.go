package responses

import (
	"ticket-booking/entities"

	"github.com/gofiber/fiber/v2"
)

type TicketResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    fiber.Map `json:"data,omitempty"`
}

func NewTicketResponse(status int, message string, tickets []*entities.Ticket, QRCode []byte) *TicketResponse {
	return &TicketResponse{
		Status:  status,
		Message: message,
		Data: fiber.Map{
			"tickets": tickets,
			"qrcode":  QRCode,
		},
	}
}
