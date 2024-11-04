package repositories

import (
	"context"
	"database/sql"
	"ticket-booking/configs/logs"
	"ticket-booking/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TicketRepository interface {
	FindAll(ctx context.Context, accountID uuid.UUID) ([]*entities.Ticket, error)
	FindByID(ctx context.Context, accountID uuid.UUID, id uint64) (*entities.Ticket, error)
	Create(ctx context.Context, ticket *entities.Ticket) error
	Validate(ctx context.Context, ticket *entities.Ticket) error
	Delete(ctx context.Context, accountID uuid.UUID, id uint64) error
}

type ticketRepository struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewTicketRepository(reader, writer *sqlx.DB) TicketRepository {
	return &ticketRepository{reader: reader, writer: writer}
}

func (t *ticketRepository) FindAll(ctx context.Context, accountID uuid.UUID) ([]*entities.Ticket, error) {
	var tickets []*entities.Ticket
	query := `SELECT * FROM tickets WHERE account_id = $1`
	if err := t.reader.SelectContext(ctx, &tickets, query, accountID); err != nil {
		logs.Error("TicketRepository.FindAll: Failed to retrieve tickets", err)
		return nil, err
	}

	return tickets, nil
}

func (t *ticketRepository) FindByID(ctx context.Context, accountID uuid.UUID, id uint64) (*entities.Ticket, error) {
	ticket := new(entities.Ticket)
	query := `SELECT * FROM tickets WHERE id = $1 AND account_id = $2`
	if err := t.reader.GetContext(ctx, ticket, query, id, accountID); err != nil {
		if err == sql.ErrNoRows {
			logs.Warn("TicketRepository.FindByID: Ticket not found")
			return nil, nil
		}
		logs.Error("TicketRepository.FindByID: Failed to retrieve ticket by ID", err)
		return nil, err
	}

	return ticket, nil
}

func (t *ticketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	query := `INSERT INTO tickets (event_id, account_id, entered, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := t.writer.QueryRowContext(ctx, query, ticket.EventID, ticket.AccountID, ticket.Entered, ticket.CreatedAt, ticket.UpdatedAt).Scan(&ticket.ID); err != nil {
		logs.Error("TicketRepository.Create: Failed to create ticket", err)
		return err
	}

	return nil
}

func (t *ticketRepository) Validate(ctx context.Context, ticket *entities.Ticket) error {
	query := `UPDATE tickets SET entered = $1, updated_at = $2 WHERE id = $3 AND account_id = $4`
	if _, err := t.writer.ExecContext(ctx, query, ticket.Entered, ticket.UpdatedAt, ticket.ID, ticket.AccountID); err != nil {
		logs.Error("TicketRepository.Validate: Failed to validate ticket", err)
		return err
	}

	return nil
}

func (t *ticketRepository) Delete(ctx context.Context, accountID uuid.UUID, id uint64) error {
	query := `DELETE FROM tickets WHERE id = $1 AND account_id = $2`
	if _, err := t.writer.ExecContext(ctx, query, id, accountID); err != nil {
		logs.Error("TicketRepository.Delete: Failed to delete ticket", err)
		return err
	}

	return nil
}
