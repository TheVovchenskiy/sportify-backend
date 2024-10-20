package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/google/uuid"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, event *models.FullEvent) error
	EditEvent(ctx context.Context, event *models.FullEvent) error
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
	result := models.NewFullEventSite(uuid.New(), request.UserID, &request.CreateEvent)

	err := a.eventStorage.CreateEvent(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("to create event: %w", err)
	}

	return result, nil
}

var ErrNotValidUser = errors.New("вы не можете изменять не свой заказ")

func (a *App) EditEventSite(ctx context.Context, request *models.RequestEventEditSite) (*models.FullEvent, error) {
	preResult := models.NewFullEventSite(request.EventID, request.UserID, &request.EventEditSite)

	eventFromDB, err := a.eventStorage.GetEvent(ctx, preResult.ID)
	if err != nil {
		return nil, fmt.Errorf("to get event: %w", err)
	}

	if eventFromDB.CreatorID != preResult.CreatorID {
		return nil, ErrNotValidUser
	}

	err = a.eventStorage.EditEvent(ctx, preResult)
	if err != nil {
		return nil, fmt.Errorf("to edit event: %w", err)
	}

	preResult.Subscribers = eventFromDB.Subscribers
	preResult.Busy = eventFromDB.Busy
	preResult.URLMessage = eventFromDB.URLMessage
	preResult.URLAuthor = eventFromDB.URLAuthor
	preResult.IsFree = eventFromDB.IsFree
	preResult.RawMessage = eventFromDB.RawMessage

	return preResult, nil
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
