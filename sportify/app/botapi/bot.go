package botapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"
)

type BotAPI struct {
	client  *http.Client
	baseURL string
	port    int
}

func NewBotAPI(baseUrl string, port int) (*BotAPI, error) {
	return &BotAPI{
		client:  http.DefaultClient,
		baseURL: baseUrl,
		port:    port,
	}, nil
}

func (api *BotAPI) EventCreated(ctx context.Context, eventCreateRequest models.EventCreatedBotRequest) error {
	reqURL := fmt.Sprintf("%s:%d/%s", api.baseURL, api.port, "event/created")

	logger, err := mylogger.Get()
	if err != nil {
		return fmt.Errorf("get logger: %w", err)
	}

	body, err := json.Marshal(eventCreateRequest)
	if err != nil {
		return fmt.Errorf("marshal event created: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	logger.WithCtx(ctx).Infow("Got response", "status", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	logger.WithCtx(ctx).Infow("Got response body", "body", string(respBody))

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("bad status code: %d", resp.StatusCode)
}

func (api *BotAPI) EventUpdated(ctx context.Context, eventUpdateRequest models.EventUpdatedBotRequest) error {
	reqURL := fmt.Sprintf("%s:%d/%s", api.baseURL, api.port, "event/updated")

	logger, err := mylogger.Get()
	if err != nil {
		return fmt.Errorf("get logger: %w", err)
	}

	body, err := json.Marshal(eventUpdateRequest)
	if err != nil {
		return fmt.Errorf("marshal event created: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	logger.WithCtx(ctx).Infow("Got response", "status", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	logger.WithCtx(ctx).Infow("Got response body", "body", string(respBody))

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("bad status code: %d", resp.StatusCode)
}

func (api *BotAPI) EventDeleted(ctx context.Context, eventDeleteRequest models.EventDeletedBotRequest) error {
	reqURL := fmt.Sprintf("%s:%d/%s", api.baseURL, api.port, "event/deleted")

	logger, err := mylogger.Get()
	if err != nil {
		return fmt.Errorf("get logger: %w", err)
	}

	body, err := json.Marshal(eventDeleteRequest)
	if err != nil {
		return fmt.Errorf("marshal event created: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	logger.WithCtx(ctx).Infow("Got response", "status", resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	logger.WithCtx(ctx).Infow("Got response body", "body", string(respBody))

	if 200 <= resp.StatusCode && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("bad status code: %d", resp.StatusCode)
}
