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
	DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error
	GetEvents(ctx context.Context) ([]models.ShortEvent, error)
	GetCreatorID(ctx context.Context, eventID uuid.UUID) (uuid.UUID, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error)
	SubscribeEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error)
}

type FileStorage interface {
	SaveFile(ctx context.Context, file []byte, fileName string) error
	Check(ctx context.Context, files []string) ([]bool, error)
}

var _ FileStorage = (*db.FileSystemStorage)(nil)

//go:generate mockgen -source=app.go -destination=mocks/app.go -package=mocks EventStorage

type App struct {
	urlPrefixFile string
	fileStorage   FileStorage
	eventStorage  EventStorage
}

//var _ EventStorage = (*db.SimpleEventStorage)(nil)

var _ EventStorage = (*db.PostgresStorage)(nil)

func NewApp(urlPrefixFile string, fileStorage FileStorage, eventStorage EventStorage) *App {
	return &App{urlPrefixFile: urlPrefixFile, eventStorage: eventStorage, fileStorage: fileStorage}
}

var (
	creatorIDTgDummy, _ = uuid.Parse("cc6edd06-43b7-4d4a-a923-dcabb819bec4")
	urlPreviewDummy     = "default_football.jpeg"
)

func (a *App) CreateEventTg(ctx context.Context, fullEvent *models.FullEvent) (*models.FullEvent, error) {
	// TODO add in db persistent map uuid to id from tg user
	fullEvent.CreatorID = creatorIDTgDummy
	fullEvent.ID = uuid.New()
	fullEvent.CreationType = models.CreationTypeTg
	// TODO try get photos from tg message and default photo to different SportType
	fullEvent.URLPreview = a.urlPrefixFile + urlPreviewDummy
	fullEvent.URLPhotos = []string{a.urlPrefixFile + urlPreviewDummy}

	err := a.eventStorage.CreateEvent(ctx, fullEvent)
	if err != nil {
		return nil, fmt.Errorf("to create event: %w", err)
	}

	return fullEvent, nil
}

func (a *App) CreateEventSite(ctx context.Context, request *models.RequestEventCreateSite) (*models.FullEvent, error) {
	result := models.NewFullEventSite(uuid.New(), request.UserID, &request.CreateEvent)

	err := a.eventStorage.CreateEvent(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("to create event: %w", err)
	}

	return result, nil
}

var ErrForbiddenEditNotYourEvent = errors.New("вы не можете изменять не свой заказ")

func (a *App) EditEventSite(ctx context.Context, request *models.RequestEventEditSite) (*models.FullEvent, error) {
	preResult := models.NewFullEventSite(request.EventID, request.UserID, &request.EventEditSite)

	eventFromDB, err := a.eventStorage.GetEvent(ctx, preResult.ID)
	if err != nil {
		return nil, fmt.Errorf("to get event: %w", err)
	}

	if eventFromDB.CreatorID != preResult.CreatorID {
		return nil, ErrForbiddenEditNotYourEvent
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

var ErrForbiddenDeleteNotYourEvent = errors.New("вы не можете удалять чужое событие")

func (a *App) DeleteEvent(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error {
	creatorID, err := a.eventStorage.GetCreatorID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("to get cretor id: %w", err)
	}

	if creatorID != userID {
		return ErrForbiddenDeleteNotYourEvent
	}

	err = a.eventStorage.DeleteEvent(ctx, userID, eventID)
	if err != nil {
		return fmt.Errorf("to delete event: %w", err)
	}

	return nil
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
