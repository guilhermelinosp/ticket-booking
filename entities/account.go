package entities

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `db:"id" json:"id" validate:"required,uuid"`
	Name      string    `db:"name" json:"name" validate:"required,min=3,max=100"`
	Email     string    `db:"email" json:"email" validate:"required,email"`
	Password  string    `db:"password" json:"password" validate:"required,min=8"`
	CreatedAt time.Time `db:"created_at" json:"created_at" validate:"required"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at" validate:"required"`
}

func NewAccount(name, email, password string) *Account {
	return &Account{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
