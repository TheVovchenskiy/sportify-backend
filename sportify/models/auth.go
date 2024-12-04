package models

import (
	"time"

	"github.com/google/uuid"
)

type UserFull struct {
	ID          uuid.UUID
	TgID        *int64
	Username    string
	Password    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FirstName   *string
	SecondName  *string
	PhotoURL    *string
	Description *string
	SportTypes  []SportType
}

func (u *UserFull) ToBotUser() *BotUser {
	return &BotUser{
		ID:       u.ID,
		Username: u.Username,
		TgID:     u.TgID,
	}
}
