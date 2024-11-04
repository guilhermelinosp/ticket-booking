package responses

import "ticket-booking/entities"

type SignUpResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Data    []*entities.Account `json:"data,omitempty"`
}

func NewSignUpResponse(status int, message string, data []*entities.Account) *SignUpResponse {
	return &SignUpResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
