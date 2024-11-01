package requests

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// EventRequest represents an event request.
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

// Validate validates the EventRequest fields.
func (e *EventRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(e)
}
