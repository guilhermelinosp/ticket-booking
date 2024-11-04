package repositories

import (
	"context"
	"ticket-booking/configs/logs"
	"ticket-booking/entities"

	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	SignUp(ctx context.Context, auth *entities.Account) (*entities.Account, error)
	FindByEmail(ctx context.Context, email string) (*entities.Account, error)
}

type accountRepository struct {
	reader *sqlx.DB
	writer *sqlx.DB
}

func NewAccountRepository(reader, writer *sqlx.DB) AccountRepository {
	return &accountRepository{reader: reader, writer: writer}
}

func (r *accountRepository) SignUp(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	query := `INSERT INTO accounts (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := r.writer.ExecContext(ctx, query, account.ID, account.Name, account.Email, account.Password, account.CreatedAt, account.UpdatedAt); err != nil {
		logs.Error("AuthRepository.SignUp: Failed to create auth", err)
		return nil, err
	}

	return account, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*entities.Account, error) {
	auth := new(entities.Account)
	query := `SELECT * FROM accounts WHERE email = $1`
	if err := r.reader.GetContext(ctx, auth, query, email); err != nil {
		logs.Error("authRepository.FindByEmail: Failed to retrieve auth by email", err)
		return nil, err
	}

	return auth, nil
}
