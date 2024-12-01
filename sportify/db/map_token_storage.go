package db

import (
	"context"
	"errors"
	"sync"
	"time"
)

type tokenInfo struct {
	username  string
	expiredAt time.Time
}

type MapTokenStorage struct {
	ttlToken time.Duration
	mu       sync.RWMutex
	tokens   map[string]tokenInfo
}

func NewMapTokenStorage(ttlToken time.Duration) *MapTokenStorage {
	return &MapTokenStorage{
		ttlToken: ttlToken,
		mu:       sync.RWMutex{},
		tokens:   make(map[string]tokenInfo),
	}
}

var (
	ErrTokenNotFound = errors.New("такой токен не найден")
	ErrTokenExpired  = errors.New("такой токен просрочен")
)

func (m *MapTokenStorage) GetUsername(_ context.Context, token string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokenInfo, ok := m.tokens[token]
	if !ok {
		return "", ErrTokenNotFound
	}

	if tokenInfo.expiredAt.Before(time.Now()) {
		delete(m.tokens, token)
		return "", ErrTokenExpired
	}

	return tokenInfo.username, nil
}

func (m *MapTokenStorage) Set(_ context.Context, token, username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tokenInfo := tokenInfo{
		username:  username,
		expiredAt: time.Now().Add(m.ttlToken),
	}

	m.tokens[token] = tokenInfo

	return nil
}
