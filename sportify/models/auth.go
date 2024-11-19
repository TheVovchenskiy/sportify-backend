package models

import (
	"github.com/google/uuid"
	"time"
)

type UserFull struct {
	ID        uuid.UUID
	TgID      *int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
