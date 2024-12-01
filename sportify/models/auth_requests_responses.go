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

type RequestLoginFromTg struct {
	TgUserID   int64  `json:"tg_user_id"`
	TgUsername string `json:"tg_username"`
	Token      string `json:"token"`
}
