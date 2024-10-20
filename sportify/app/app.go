package app

import (
	"context"
	"fmt"

	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/google/uuid"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, event *models.FullEvent) error
	GetEvents(ctx context.Context) ([]models.ShortEvent, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error)
	SubscribeEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error)
}

//go:generate mockgen -source=app.go -destination=mocks/app.go -package=mocks EventStorage

type App struct {
	eventStorage EventStorage
}

//var _ EventStorage = (*db.SimpleEventStorage)(nil)

var _ EventStorage = (*db.PostgresStorage)(nil)

func NewApp(eventStorage EventStorage) *App {
	return &App{eventStorage: eventStorage}
}

func (a *App) CreateEventSite(ctx context.Context, request *models.RequestEventCreateSite) (*models.FullEvent, error) {
	result := models.FullEvent{
		ShortEvent: models.ShortEvent{
			ID:          uuid.New(),
			CreatorID:   request.UserID,
			SportType:   request.CreateEvent.SportType,
			Address:     request.CreateEvent.Address,
			Date:        request.CreateEvent.Date,
			StartTime:   request.CreateEvent.StartTime,
			EndTime:     request.CreateEvent.EndTime,
			Price:       request.CreateEvent.Price,
			IsFree:      models.IsFreePrice(request.CreateEvent.Price),
			GameLevel:   request.CreateEvent.GameLevel,
			Capacity:    request.CreateEvent.Capacity,
			Busy:        0,
			Subscribers: make([]uuid.UUID, 0),
			URLPreview:  request.CreateEvent.URLPreview,
			URLPhotos:   request.CreateEvent.URLPhotos,
		},
		CreationType: models.CreationTypeSite,
		Description:  request.CreateEvent.Description,
	}

	err := a.eventStorage.CreateEvent(ctx, &result)
	if err != nil {
		return nil, fmt.Errorf("to create event: %w", err)
	}

	return &result, nil
}

func (a *App) GetEvents(ctx context.Context) ([]models.ShortEvent, error) {
	return a.eventStorage.GetEvents(ctx)
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error) {
	return a.eventStorage.GetEvent(ctx, id)
}

func (a *App) SubscribeEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error) {
	return a.eventStorage.SubscribeEvent(ctx, id, userID, subscribe)
}
