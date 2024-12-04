package models

import (
	"time"

	"github.com/google/uuid"
)

type RequestLogin struct {
	Username    string `json:"user"`
	PasswordRaw string `json:"passwd"`
}

type ResponseSuccessLogin struct {
	UserID   uuid.UUID `json:"user_id"`
	TgUserID *int64    `json:"tg_user_id"`
	Username string    `json:"username"`
}

type TgUpdate struct {
	UpdateID int "json:\"update_id\""
	Message  struct {
		Chat struct {
			ID   int    "json:\"id\""
			Name string "json:\"first_name\""
			Type string "json:\"type\""
		} "json:\"chat\""
		Text string "json:\"text\""
	} "json:\"message\""
}

type TgUpdateWrapper struct {
	TgUpdate
	ExpirationTime time.Time
}

type TgRequestAuth struct {
	TgUpdate   TgUpdate `json:"tg_update"`
	TgUserID   int64    `json:"tg_user_id"`
	TgUsername string   `json:"tg_username"`
	Token      string   `json:"token"`
}

// type TgResponseAuth struct {
// 	UserID   uuid.UUID `json:"user_id"`
// 	TgUserID int64     `json:"tg_user_id"`
// 	Username string    `json:"username"`
// }

// type RequestLoginFromTg struct {
// 	TgUserID   int64  `json:"tg_user_id"`
// 	TgUsername string `json:"tg_username"`
// 	Token      string `json:"token"`
// }
