package responses

import (
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
}

// NewEventResponse creates a new EventResponse with the provided details.
func NewEventResponse(id uuid.UUID, title, location string, date time.Time) *EventResponse {
	return &EventResponse{
		ID:       id,
		Title:    title,
		Date:     date,
		Location: location,
	}
}
