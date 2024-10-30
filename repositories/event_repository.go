package repositories

import (
	"context"
	"ticket-booking/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewEventRepository(reader, writer *sqlx.DB) *EventRepository {
	return &EventRepository{
		reader: reader,
		writer: writer,
	}
}

// FindAll implements interfaces.EventRepository.
func (r *EventRepository) FindAll(ctx context.Context) ([]*entities.Event, error) {
	events := []*entities.Event{}
	err := r.reader.SelectContext(ctx, &events, "SELECT * FROM events")
	if err != nil {
		return nil, err
	}
	return events, nil
}

// FindByID implements interfaces.EventRepository.
func (r *EventRepository) FindByID(ctx context.Context, eventId uuid.UUID) (*entities.Event, error) {
	event := &entities.Event{}
	err := r.reader.GetContext(ctx, event, "SELECT * FROM events WHERE id = ?", eventId)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Create implements interfaces.EventRepository.
func (r *EventRepository) Create(ctx context.Context, event *entities.Event) (entities.Event, error) {
	query := "INSERT INTO events (id, title, location, date, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.writer.ExecContext(ctx, query, event.ID, event.Title, event.Location, event.Date, event.CreatedAt, event.UpdatedAt)
	if err != nil {
		return entities.Event{}, err
	}
	return *event, nil
}

// Update implements interfaces.EventRepository.
func (r *EventRepository) Update(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	query := "UPDATE events SET title = ?, location = ?, date = ?, updated_at = ? WHERE id = ?"
	_, err := r.writer.ExecContext(ctx, query, event.Title, event.Location, event.Date, event.UpdatedAt, event.ID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Delete implements interfaces.EventRepository.
func (r *EventRepository) Delete(ctx context.Context, eventId uuid.UUID) error {
	_, err := r.writer.ExecContext(ctx, "DELETE FROM events WHERE id = ?", eventId)
	return err
}
