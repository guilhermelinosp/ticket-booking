package repositories

import (
	"context"
	"ticket-booking/configs/logs"
	"ticket-booking/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// EventRepository defines methods for accessing events data.
type EventRepository interface {
	FindAll(ctx context.Context) ([]*entities.Event, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Event, error)
	Create(ctx context.Context, event *entities.Event) (*entities.Event, error)
	Update(ctx context.Context, event *entities.Event) (*entities.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// eventRepository implements EventRepository.
type eventRepository struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

// NewEventRepository initializes a new repository for events.
func NewEventRepository(reader, writer *sqlx.DB) EventRepository {
	return &eventRepository{reader: reader, writer: writer}
}

func (r *eventRepository) FindAll(ctx context.Context) ([]*entities.Event, error) {
	var events []*entities.Event
	if err := r.reader.SelectContext(ctx, &events, "SELECT * FROM events"); err != nil {
		logs.Error("EventRepository.FindAll: Failed to retrieve events", err)
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Event, error) {
	event := new(entities.Event)
	if err := r.reader.GetContext(ctx, event, "SELECT * FROM events WHERE id = $1", id); err != nil {
		logs.Error("EventRepository.FindByID: Failed to retrieve event by ID", err)
		return nil, err
	}

	return event, nil
}

func (r *eventRepository) Create(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	query := `INSERT INTO events (id, title, location, date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := r.writer.ExecContext(ctx, query, event.ID, event.Title, event.Location, event.Date, event.CreatedAt, event.UpdatedAt); err != nil {
		logs.Error("EventRepository.Create: Failed to create event", err)
		return nil, err
	}

	return event, nil
}

func (r *eventRepository) Update(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	query := `UPDATE events SET title = $1, location = $2, date = $3, updated_at = $4 WHERE id = $5`
	if _, err := r.writer.ExecContext(ctx, query, event.Title, event.Location, event.Date, event.UpdatedAt, event.ID); err != nil {
		logs.Error("EventRepository.Update: Failed to update event", err)
		return nil, err
	}

	return event, nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.writer.ExecContext(ctx, "DELETE FROM events WHERE id = $1", id); err != nil {
		logs.Error("EventRepository.Delete: Failed to delete event", err)
		return err
	}

	return nil
}
