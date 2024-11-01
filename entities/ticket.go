package entities

import (
	"time"

	"github.com/google/uuid"
)

// Ticket represents a ticket.
type Ticket struct {
	ID        uuid.UUID `db:"id" json:"id" valid:"uuid"`
	EventID   uuid.UUID `db:"event_id" json:"event_id" valid:"uuid"`
	Event     *Event    `db:"event" json:"event" valid:"-" relation:"event_id" fk:"id" `
	Entered   bool      `db:"entered" json:"entered" valid:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at" valid:"required"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" valid:"required"`
}

// NewTicket creates a new instance of Ticket.
func NewTicket(eventID uuid.UUID) *Ticket {
	return &Ticket{
		ID:        uuid.New(),
		EventID:   eventID,
		Entered:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
