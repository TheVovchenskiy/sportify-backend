package yookassa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheVovchenskiy/sportify-backend/models"
)

type Client struct {
	shopID       string
	agentID      string
	tokenPayment string
	tokenPayout  string
	httpClient   *http.Client
}

func NewClient(shopID, agentID, tokenPayment, tokenPayout string) *Client {
	return &Client{
		shopID:       shopID,
		agentID:      agentID,
		tokenPayment: tokenPayment,
		tokenPayout:  tokenPayout,
		httpClient:   http.DefaultClient,
	}
}

const urlPayment = "https://api.yookassa.ru/v3/payments"

//nolint:err113
func (c *Client) DoPayment(
	ctx context.Context,
	idempotencyKey,
	redirectURL string,
	amount float64,
) (*models.Payment, error) {
	paymentRequest := NewRequestPayment(redirectURL, amount)

	payload, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPayment, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.shopID, c.tokenPayment)
	req.Header.Set("Idempotence-Key", idempotencyKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("to do request: %w", err)
	}
	defer resp.Body.Close()

	//body, err := io.ReadAll(resp.Body)
	//fmt.Println(string(body), err)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var responsePayment ResponsePayment
	if err := json.NewDecoder(resp.Body).Decode(&responsePayment); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	responseAmount, err := strconv.ParseFloat(responsePayment.Amount.Value, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	if int64(responseAmount) != int64(amount) {
		return nil, fmt.Errorf("unexpected amount: %f", responseAmount)
	}

	if responsePayment.Status == "succeeded" {
		responsePayment.Status = "paid"
	} else if responsePayment.Status == "canceled" {
		responsePayment.Status = "canceled"
	} else {
		responsePayment.Status = "pending"
	}

	return &models.Payment{ //nolint:exhaustruct
		ID:              responsePayment.ID,
		ConfirmationURL: responsePayment.Confirmation.ConfirmationURL,
		Status:          models.PaymentStatus(responsePayment.Status), //TODO will be good check in yookassa is right only our statuses?
		Amount:          int64(responseAmount),
	}, nil
}
