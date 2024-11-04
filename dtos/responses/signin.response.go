package responses

import "time"

// SignInResponse represents the response for a sign-in request.
type SignInResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    []*TokenResponse `json:"data,omitempty"`
}

// TokenResponse represents the response containing the token and refresh token.
type TokenResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh"`
	ExpiryDate   time.Time `json:"expiry"`
}

// NewTokenResponse creates a new instance of TokenResponse.
func NewTokenResponse(token, refresh string, expiry time.Time) *TokenResponse {
	return &TokenResponse{
		Token:        token,
		RefreshToken: refresh,
		ExpiryDate:   expiry,
	}
}

// NewSignInResponse creates a new instance of SignInResponse.
func NewSignInResponse(status int, message string, tokenResponse []*TokenResponse) *SignInResponse {
	return &SignInResponse{
		Status:  status,
		Message: message,
		Data:    tokenResponse,
	}
}
