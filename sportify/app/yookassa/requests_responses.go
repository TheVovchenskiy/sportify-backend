package yookassa

import (
	"fmt"

	"github.com/google/uuid"
)

type RequestPayment struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Capture      bool `json:"capture"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
}

func NewRequestPayment(redirectURL string, amount float64) *RequestPayment {
	return &RequestPayment{
		Amount: struct {
			Value    string `json:"value"`
			Currency string `json:"currency"`
		}{
			Value:    fmt.Sprintf("%.2f", amount),
			Currency: "RUB",
		},
		Capture: true,
		Confirmation: struct {
			Type      string `json:"type"`
			ReturnURL string `json:"return_url"`
		}{
			Type:      "redirect",
			ReturnURL: redirectURL,
		},
	}
}

type ResponsePayment struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}
