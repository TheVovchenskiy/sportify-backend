package models

import (
	"time"

	"github.com/google/uuid"
)

type UserFull struct {
	ID        uuid.UUID
	TgID      *int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *UserFull) ToBotUser() *BotUser {
	return &BotUser{
		ID:       u.ID,
		Username: u.Username,
		TgID:     u.TgID,
	}
}
