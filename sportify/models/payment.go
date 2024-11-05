package models

import "github.com/google/uuid"

type PaymentStatus string

const (
	PaymentStatusPaid      = "paid"
	PaymentStatusCancelled = "cancelled"
	PaymentStatusPending   = "pending"
)

type Payment struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	EventID         uuid.UUID     `json:"event_id"`
	ConfirmationURL string        `json:"confirmation_url"`
	Status          PaymentStatus `json:"status"`
	Amount          int64         `json:"amount"`
}
