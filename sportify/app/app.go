package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/app/botapi"
	"github.com/TheVovchenskiy/sportify-backend/app/yookassa"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	"github.com/google/uuid"
)

type EventStorage interface {
	CreateEvent(ctx context.Context, event *models.FullEvent) error
	EditEvent(ctx context.Context, event *models.FullEvent) error
	DeleteEvent(ctx context.Context, userID, eventID uuid.UUID) error
	GetEvents(ctx context.Context) ([]models.ShortEvent, error)
	GetCreatorID(ctx context.Context, eventID uuid.UUID) (uuid.UUID, error)
	FindEvents(ctx context.Context, filterParams *models.FilterParams) ([]models.ShortEvent, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error)
	SubscribeEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error)
	AddUserPaid(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	SetCoordinates(ctx context.Context, latitude, longitude string, id uuid.UUID) error
}

//var _ EventStorage = (*db.SimpleEventStorage)(nil)

var _ EventStorage = (*db.PostgresStorage)(nil)

type FileStorage interface {
	SaveFile(ctx context.Context, file []byte, fileName string) error
	Check(ctx context.Context, files []string) ([]bool, error)
}

var _ FileStorage = (*db.FileSystemStorage)(nil)

type BotAPI interface {
	EventCreated(ctx context.Context, eventCreateRequest models.EventCreatedBotRequest) error
}

var _ BotAPI = (*botapi.BotAPI)(nil)

type YookassaClient interface {
	DoPayment(ctx context.Context, idempotencyKey, redirectURL string, amount float64) (*models.Payment, error)
}

var _ YookassaClient = (*yookassa.Client)(nil)

type PaymentPayoutStorage interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPayment(ctx context.Context, id uuid.UUID) (*models.Payment, error)
	UpdateStatusPayment(ctx context.Context, id uuid.UUID, status models.PaymentStatus) error
}

var _ PaymentPayoutStorage = (*db.PostgresPaymentPayoutStorage)(nil)

//go:generate mockgen -source=app.go -destination=mocks/app.go -package=mocks EventStorage

type App struct {
	urlPrefixFile        string
	fileStorage          FileStorage
	eventStorage         EventStorage
	paymentPayoutStorage PaymentPayoutStorage
	yookassaClient       YookassaClient
	httpClient           *http.Client
	logger               *mylogger.MyLogger
	botAPI               BotAPI
}

func NewApp(
	urlPrefixFile string,
	fileStorage FileStorage,
	eventStorage EventStorage,
	logger *mylogger.MyLogger,
	botAPI BotAPI,
	// paymentPayoutStorage PaymentPayoutStorage,
	// yookassaClient YookassaClient,
) *App {
	app := &App{
		urlPrefixFile: urlPrefixFile,
		eventStorage:  eventStorage,
		fileStorage:   fileStorage,
		httpClient:    http.DefaultClient,
		logger:        logger,
		botAPI:        botAPI,
		// paymentPayoutStorage: paymentPayoutStorage,
		// yookassaClient:       yookassaClient,
	}

	// TODO add context to cancel
	go func() {
		defer func() {
			if pan := recover(); pan != nil {
				logger.Errorf("panic: %v", pan)
			}
		}()
		app.RefreshCoordinates(context.TODO(), time.Second*60)
	}()

	return app
}

var (
	creatorIDTgDummy, _ = uuid.Parse("cc6edd06-43b7-4d4a-a923-dcabb819bec4")
	urlPreviewDummy     = "default_football.jpeg"
)

func (a *App) CreateEventTg(ctx context.Context, fullEvent *models.FullEvent) (*models.FullEvent, error) {
	// TODO add in db persistent map uuid to id from tg user
	fullEvent.CreatorID = creatorIDTgDummy
	fullEvent.ID = uuid.New()
	// TODO try get photos from tg message and default photo to different SportType
	fullEvent.CreationType = models.CreationTypeTg
	fullEvent.IsFree = models.IsFreePrice(fullEvent.Price)
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

	if request.Tg != nil {
		chatID, err := strconv.ParseInt(*request.Tg.ChatID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("to parse chatID: %w", err)
		}

		eventCreateRequest := models.EventCreatedBotRequest{
			TgChatID: &chatID,
			TgUserID: request.Tg.UserID,
			Event:    result.ShortEvent,
		}

		err = a.botAPI.EventCreated(ctx, eventCreateRequest)
		if err != nil {
			// TODO: maybe we should consider not to send 500 error to user in this case
			return nil, fmt.Errorf("to send to bot created event: %w", err)
		}
	}

	return result, nil
}

var ErrForbiddenEditNotYourEvent = errors.New("вы не можете изменять не свое событие")

func (a *App) EditEventSite(ctx context.Context, request *models.RequestEventEditSite) (*models.FullEvent, error) {
	eventFromDB, err := a.eventStorage.GetEvent(ctx, request.EventID)
	if err != nil {
		return nil, fmt.Errorf("to get event: %w", err)
	}

	if eventFromDB.CreatorID != request.UserID {
		return nil, ErrForbiddenEditNotYourEvent
	}

	if len(request.EventEditSite.GameLevels) == 0 {
		request.EventEditSite.GameLevels = eventFromDB.GameLevels
	}

	if len(request.EventEditSite.URLPhotos) == 0 {
		request.EventEditSite.URLPhotos = eventFromDB.URLPhotos
	}

	preResult := &models.FullEvent{
		ShortEvent: models.ShortEvent{
			ID:          request.EventID,
			CreatorID:   eventFromDB.CreatorID,
			SportType:   common.NewValWithFallback(request.EventEditSite.SportType, &eventFromDB.SportType),
			Address:     common.NewValWithFallback(request.EventEditSite.Address, &eventFromDB.Address),
			DateAndTime: common.NewValWithFallback(request.EventEditSite.DateAndTime, &eventFromDB.DateAndTime),
			Price:       common.Ref(common.NewValWithFallback(request.EventEditSite.Price, eventFromDB.Price)),
			GameLevels:  request.EventEditSite.GameLevels,
			Capacity:    common.Ref(common.NewValWithFallback(request.EventEditSite.Capacity, eventFromDB.Capacity)),
			URLPreview:  common.NewValWithFallback(request.EventEditSite.URLPreview, &eventFromDB.URLPreview),
			URLPhotos:   request.EventEditSite.URLPhotos,
		},
		CreationType: eventFromDB.CreationType,
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

func (a *App) FindEvents(ctx context.Context, filterParams *models.FilterParams) ([]models.ShortEvent, error) {
	return a.eventStorage.FindEvents(ctx, filterParams)
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error) {
	return a.eventStorage.GetEvent(ctx, id)
}

func (a *App) SubscribeEvent(ctx context.Context, id uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error) {
	return a.eventStorage.SubscribeEvent(ctx, id, userID, subscribe)
}

var ErrPayFree = errors.New("вы не можете оплатить бесплатное событие")

func (a *App) PayEvent(ctx context.Context, request *models.RequestEventPay) (*models.ResponseEventPay, error) {
	fullEvent, err := a.GetEvent(ctx, request.EventID)
	if err != nil {
		return nil, fmt.Errorf("to get event: %w", err)
	}

	if fullEvent.IsFree || fullEvent.Price == nil {
		return nil, ErrPayFree
	}

	amount := float64(*fullEvent.Price)

	idempotencyKey := []byte(request.EventID.String() + request.UserID.String())[:64]

	payment, err := a.yookassaClient.DoPayment(ctx, string(idempotencyKey), request.RedirectURL, amount)
	if err != nil {
		return nil, fmt.Errorf("to do payment: %w", err)
	}

	payment.UserID = request.UserID
	payment.EventID = request.EventID

	err = a.paymentPayoutStorage.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("to create payment: %w", err)
	}

	return &models.ResponseEventPay{
		ID:              payment.ID,
		ConfirmationURL: payment.ConfirmationURL,
	}, nil
}

var (
	mapPayment = map[string]struct{}{}
	muPayment  = sync.RWMutex{}
)

func (a *App) GetPayment(ctx context.Context, paymentID uuid.UUID) (*models.ResponsesPayment, error) {
	go func() {
		// TODO add check payment status, do payout
		muPayment.Lock()
		if _, ok := mapPayment[paymentID.String()]; ok {
			muPayment.Unlock()
			return
		}

		mapPayment[paymentID.String()] = struct{}{}
		go func() {
			time.Sleep(time.Second * 30)
			ctx := context.TODO()
			err := a.paymentPayoutStorage.UpdateStatusPayment(ctx, paymentID, models.PaymentStatusPaid)
			if err != nil {
				fmt.Println("update status payment: ", err)
			}

			payment, err := a.paymentPayoutStorage.GetPayment(ctx, paymentID)
			if err != nil {
				fmt.Println("get payment: ", err)
			}

			err = a.eventStorage.AddUserPaid(ctx, payment.EventID, payment.UserID)
			if err != nil {
				fmt.Println("add user paid: ", err)
			}
		}()
		muPayment.Unlock()
	}()

	payment, err := a.paymentPayoutStorage.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("to get payment: %w", err)
	}

	return &models.ResponsesPayment{PaymentStatus: payment.Status}, nil
}
