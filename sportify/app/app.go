package app

import (
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/google/uuid"
)

type EventStorage interface {
	AddEvent(event models.FullEvent) error
	GetEvents() ([]models.ShortEvent, error)
	GetEvent(id uuid.UUID) (*models.FullEvent, error)
	SubscribeEvent(id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error)
}

type App struct {
	eventStorage EventStorage
}

var _ EventStorage = (*db.SimpleEventStorage)(nil)

func NewApp(eventStorage EventStorage) *App {
	return &App{eventStorage: eventStorage}
}

func (a *App) AddEvent(event models.FullEvent) error {
	return a.eventStorage.AddEvent(event)
}

func (a *App) GetEvents() ([]models.ShortEvent, error) {
	return a.eventStorage.GetEvents()
}

func (a *App) GetEvent(id uuid.UUID) (*models.FullEvent, error) {
	return a.eventStorage.GetEvent(id)
}

func (a *App) SubscribeEvent(id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error) {
	return a.eventStorage.SubscribeEvent(id, userID, subscribe)
}
