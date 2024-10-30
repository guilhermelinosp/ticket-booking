package interfaces

import (
	"context"
	"ticket-booking/entities"

	"github.com/google/uuid"
)

type EventRepository interface {
	FindAll(ctx context.Context) ([]*entities.Event, error)
	FindByID(ctx context.Context, eventId uuid.UUID) (*entities.Event, error)
	Create(ctx context.Context, event *entities.Event) (entities.Event, error)
	Update(ctx context.Context, event *entities.Event) (*entities.Event, error)
	Delete(ctx context.Context, eventId uuid.UUID) error
}

