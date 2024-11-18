package models

import "github.com/google/uuid"

type RequestLogin struct {
	Username    string `json:"user"`
	PasswordRaw string `json:"passwd"`
}

type ResponseSuccessLogin struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
