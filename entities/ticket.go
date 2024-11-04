package entities

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID        uint64    `db:"id" json:"id" valid:"uint"`
	EventID   uint64    `db:"event_id" json:"event_id" valid:"uint" relation:"event_id" fk:"id"`
	Event     *Event    `db:"event" json:"event" valid:"-" relation:"event_id" fk:"id" `
	AccountID uuid.UUID `db:"account_id" json:"account_id" valid:"uuid" relation:"account_id" fk:"id"`
	Entered   bool      `db:"entered" json:"entered" valid:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at" valid:"required"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" valid:"required"`
}

func NewTicket(eventID uint64, accountID uuid.UUID) *Ticket {
	return &Ticket{
		EventID:   eventID,
		AccountID: accountID,
		Entered:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
