package telegramapi

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/timeout"
	"github.com/go-pkgz/auth/provider"
)

type TelegramAPIDummy struct {
	// provider.TelegramAPI
	updatesStorage []*models.TgUpdateWrapper
	muUpdates      sync.Mutex

	results   map[int]error
	muResults sync.Mutex
}

var _ provider.TelegramAPI = (*TelegramAPIDummy)(nil)

func NewTelegramAPIDummy() *TelegramAPIDummy {
	return &TelegramAPIDummy{
		updatesStorage: make([]*models.TgUpdateWrapper, 0),
		muUpdates:      sync.Mutex{},

		results:   make(map[int]error),
		muResults: sync.Mutex{},
	}
}

func (t *TelegramAPIDummy) AddUpdate(update *models.TgUpdateWrapper) {
	t.muUpdates.Lock()
	defer t.muUpdates.Unlock()

	t.updatesStorage = append(t.updatesStorage, update)
}

var ErrNoResult = errors.New("no result")

func (t *TelegramAPIDummy) getResult(chatID int) error {
	t.muResults.Lock()
	defer t.muResults.Unlock()

	result, ok := t.results[chatID]
	if !ok {
		return ErrNoResult
	}

	delete(t.results, chatID)

	return result
}

var (
	getResultLoopTicker  = 100 * time.Millisecond
	getResultLoopTimeout = 10 * time.Second
)

func (t *TelegramAPIDummy) GetResult(chatID int) error {
	getResultLoop := func() error {
		ticker := time.NewTicker(getResultLoopTicker)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				result := t.getResult(chatID)
				if !errors.Is(result, ErrNoResult) {
					return result
				}
			}
		}
	}

	return timeout.Timeout(getResultLoopTimeout, getResultLoop)
}

func (t *TelegramAPIDummy) GetUpdates(_ context.Context) (*provider.TelegramUpdate, error) {
	t.muUpdates.Lock()
	defer t.muUpdates.Unlock()

	result := provider.TelegramUpdate{
		Result: []struct {
			UpdateID int "json:\"update_id\""
			Message  struct {
				Chat struct {
					ID   int    "json:\"id\""
					Name string "json:\"first_name\""
					Type string "json:\"type\""
				} "json:\"chat\""
				Text string "json:\"text\""
			} "json:\"message\""
		}{},
	}

	for _, update := range t.updatesStorage {
		if update.ExpirationTime.Before(time.Now()) {
			continue
		}

		result.Result = append(result.Result, update.TgUpdate)
	}

	t.updatesStorage = make([]*models.TgUpdateWrapper, 0)

	return &result, nil
}

func (t *TelegramAPIDummy) Avatar(_ context.Context, _ int) (string, error) {
	return "", nil
}

func (t *TelegramAPIDummy) Send(_ context.Context, chatID int, text string) error {
	t.muResults.Lock()
	defer t.muResults.Unlock()

	if text == "" {
		t.results[chatID] = nil
		return nil
	}

	t.results[chatID] = errors.New(text)

	return nil
}

func (t *TelegramAPIDummy) BotInfo(_ context.Context) (*provider.BotInfo, error) {
	return &provider.BotInfo{
		Username: "ond_sportify_bot", // TODO: handle in config
	}, nil
}
