package requests

import (
	"time"
)

type EventRequest struct {
	Title    string    `json:"title" validate:"required,min=3,max=100"`
	Date     time.Time `json:"date" validate:"required"`
	Location string    `json:"location" validate:"required,min=3,max=100"`
}

// NewEventRequest creates a new instance of EventRequest
func NewEventRequest(title, location string, date time.Time) *EventRequest {
	return &EventRequest{
		Title:    title,
		Location: location,
		Date:     date,
	}
}
