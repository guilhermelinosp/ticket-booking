package requests

import "github.com/go-playground/validator/v10"

// SignInRequest represents a signin request.
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// NewSignInRequest creates a new instance of SignInRequest.
func NewSignInRequest(email, password string) *SignInRequest {
	return &SignInRequest{
		Email:    email,
		Password: password,
	}
}

// Validate validates the SignUpRequest fields.
func (s *SignInRequest) Validate() error {
	return validator.New().Struct(s)
}
