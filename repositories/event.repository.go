package repositories

import (
	"context"
	"database/sql"
	"ticket-booking/configs/logs"
	"ticket-booking/entities"

	"github.com/jmoiron/sqlx"
)

type EventRepository interface {
	FindAll(ctx context.Context) ([]*entities.Event, error)
	FindByID(ctx context.Context, id uint64) (*entities.Event, error)
	Create(ctx context.Context, event *entities.Event) error
	Update(ctx context.Context, event *entities.Event) error
	Delete(ctx context.Context, id uint64) error
}

type eventRepository struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewEventRepository(reader, writer *sqlx.DB) EventRepository {
	return &eventRepository{reader: reader, writer: writer}
}

func (r *eventRepository) FindAll(ctx context.Context) ([]*entities.Event, error) {
	var events []*entities.Event
	query := `SELECT * FROM events`
	if err := r.reader.SelectContext(ctx, &events, query); err != nil {
		logs.Error("EventRepository.FindAll: Failed to retrieve events", err)
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) FindByID(ctx context.Context, id uint64) (*entities.Event, error) {
	event := new(entities.Event)
	query := `SELECT * FROM events WHERE id = $1`
	if err := r.reader.GetContext(ctx, event, query, id); err != nil {
		if err == sql.ErrNoRows {
			logs.Warn("EventRepository.FindByID: Event not found")
			return nil, nil
		}
		logs.Error("EventRepository.FindByID: Failed to retrieve event by ID", err)
		return nil, err
	}

	return event, nil
}

func (r *eventRepository) Create(ctx context.Context, event *entities.Event) error {
	query := `INSERT INTO events (title, location, date, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := r.writer.QueryRowContext(ctx, query, event.Title, event.Location, event.Date, event.CreatedAt, event.UpdatedAt).Scan(&event.ID); err != nil {
		logs.Error("EventRepository.Create: Failed to create event", err)
		return err
	}

	return nil
}

func (r *eventRepository) Update(ctx context.Context, event *entities.Event) error {
	query := `UPDATE events SET title = $1, location = $2, date = $3, updated_at = $4 WHERE id = $5`
	if _, err := r.writer.ExecContext(ctx, query, event.Title, event.Location, event.Date, event.UpdatedAt, event.ID); err != nil {
		logs.Error("EventRepository.Update: Failed to update event", err)
		return err
	}

	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM events WHERE id = $1`
	if _, err := r.writer.ExecContext(ctx, query, id); err != nil {
		logs.Error("EventRepository.Delete: Failed to delete event", err)
		return err
	}

	return nil
}
