package entities

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID `db:"id" json:"id" valid:"uuid"`
	Title     string    `db:"title" json:"title" valid:"string,required"`
	Location  string    `db:"location" json:"location" valid:"string,required"`
	Date      time.Time `db:"date" json:"date" valid:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at" valid:"required"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" valid:"required"`
}

func NewEvent(title, location string, date time.Time) *Event {
	return &Event{
		ID:        uuid.New(),
		Title:     title,
		Location:  location,
		Date:      date,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
