package requests

import "github.com/go-playground/validator/v10"

// SignupRequest represents a signup request.
type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// NewSignUpRequest creates a new instance of SignUpRequest.
func NewSignUpRequest(name, email, password string) *SignUpRequest {
	return &SignUpRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// Validate validates the SignUpRequest fields.
func (s *SignUpRequest) Validate() error {
	return validator.New().Struct(s)
}
