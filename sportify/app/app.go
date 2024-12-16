package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/app/botapi"
	"github.com/TheVovchenskiy/sportify-backend/app/yookassa"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
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
	EventUpdated(ctx context.Context, eventUpdateRequest models.EventUpdatedBotRequest) error
	EventDeleted(ctx context.Context, eventDeleteRequest models.EventDeletedBotRequest) error
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

type TokenStorage interface {
	GetUsername(ctx context.Context, token string) (string, error)
	Set(ctx context.Context, token, username string) error
}

var _ TokenStorage = (*db.MapTokenStorage)(nil)

//go:generate mockgen -source=app.go -destination=mocks/app.go -package=mocks EventStorage

type App struct {
	urlPrefixFile        string
	fileStorage          FileStorage
	eventStorage         EventStorage
	paymentPayoutStorage PaymentPayoutStorage
	authStorage          AuthStorage
	tokenStorage         TokenStorage
	yookassaClient       YookassaClient
	httpClient           *http.Client
	logger               *mylogger.MyLogger
	muFindByAddress      *sync.Mutex
	botAPI               BotAPI
}

func NewApp(
	urlPrefixFile string,
	fileStorage FileStorage,
	eventStorage EventStorage,
	authStorage AuthStorage,
	tokenStorage TokenStorage,
	logger *mylogger.MyLogger,
	botAPI BotAPI,
	// paymentPayoutStorage PaymentPayoutStorage,
	// yookassaClient YookassaClient,
) *App {
	app := &App{
		urlPrefixFile:   urlPrefixFile,
		eventStorage:    eventStorage,
		fileStorage:     fileStorage,
		authStorage:     authStorage,
		httpClient:      http.DefaultClient,
		tokenStorage:    tokenStorage,
		logger:          logger,
		muFindByAddress: &sync.Mutex{},
		botAPI:          botAPI,
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

func (a *App) getBotUser(ctx context.Context, userID uuid.UUID) (*models.BotUser, error) {
	fullUser, err := a.authStorage.GetUserFullByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("to get user: %w", err)
	}

	return fullUser.ToBotUser(), nil
}

func (a *App) getBotEvent(ctx context.Context, eventID uuid.UUID) (*models.BotEvent, error) {
	fullEvent, err := a.GetEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("to get event: %w", err)
	}

	hashtags := GenerateHashtags(&fullEvent.ShortEvent)

	creator, err := a.getBotUser(ctx, fullEvent.CreatorID)
	if err != nil {
		return nil, fmt.Errorf("to get creator: %w", err)
	}

	subscribers := []*models.BotUser{}
	// TODO: make batch request
	for _, subscriberID := range fullEvent.Subscribers {
		subscriber, err := a.getBotUser(ctx, subscriberID)
		if err != nil {
			return nil, fmt.Errorf("to get subscriber: %w", err)
		}

		subscribers = append(subscribers, subscriber)
	}

	return fullEvent.ToBotEvent(creator, subscribers, &hashtags), nil
}

func (a *App) onEventUpdate(ctx context.Context, eventID uuid.UUID) error {
	botEvent, err := a.getBotEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("to get bot event: %w", err)
	}

	eventUpdated := models.EventUpdatedBotRequest{
		Event: *botEvent,
	}

	err = a.botAPI.EventUpdated(ctx, eventUpdated)
	if err != nil {
		return fmt.Errorf("to send event updated: %w", err)
	}

	return nil
}

func (a *App) onEventDelete(ctx context.Context, eventID uuid.UUID) error {
	eventDeleted := models.EventDeletedBotRequest{
		EventID: eventID,
	}

	err := a.botAPI.EventDeleted(ctx, eventDeleted)
	if err != nil {
		return fmt.Errorf("to send event deleted: %w", err)
	}

	return nil
}

func (a *App) getDefaultEventPhoto(sportType models.SportType) string {
	result := a.urlPrefixFile
	switch sportType {
	case models.SportTypeFootball:
		return result + "default_football.jpeg"
	case models.SportTypeBasketball:
		return result + "default_basketball.png"
	case models.SportTypeVolleyball:
		return result + "default_volleyball.jpg"
	}

	return result + "default_football.jpeg"
}

func (a *App) CreateEventSite(ctx context.Context, request *models.RequestEventCreateSite) (*models.FullEvent, error) {
	result := models.NewFullEventSite(uuid.New(), request.UserID, &request.CreateEvent)

	if result.URLPreview == "" || len(result.URLPhotos) == 0 {
		defaultPhoto := a.getDefaultEventPhoto(result.SportType)
		result.URLPreview = defaultPhoto
		result.URLPhotos = []string{defaultPhoto}
	}

	err := a.eventStorage.CreateEvent(ctx, result)
	if err != nil {
		return nil, fmt.Errorf("to create event: %w", err)
	}

	if request.Tg != nil {
		a.logger.WithCtx(ctx).Infow("Creating tg event", "event_id", result.ID, "error", err)

		chatID, err := strconv.ParseInt(*request.Tg.ChatID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("to parse chatID: %w", err)
		}

		botEvent, err := a.getBotEvent(ctx, result.ID)
		if err != nil {
			a.logger.WithCtx(ctx).Warnw("Unable to get bot event", "event_id", result.ID, "error", err)
			return result, nil
		}

		eventCreateRequest := models.EventCreatedBotRequest{
			TgChatID: &chatID,
			// TgUserID: request.Tg.UserID,
			Event: *botEvent,
		}

		err = a.botAPI.EventCreated(ctx, eventCreateRequest)
		if err != nil {
			a.logger.WithCtx(ctx).Warnw("Unable to create bot event", "event_id", result.ID, "error", err)
			return result, nil
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

	if request.EventEditSite.SportType == nil {
		request.EventEditSite.SportType = &eventFromDB.SportType
	}

	if request.EventEditSite.URLPreview == nil {
		request.EventEditSite.URLPreview = &eventFromDB.URLPreview
	}

	if len(request.EventEditSite.URLPhotos) == 0 {
		request.EventEditSite.URLPhotos = eventFromDB.URLPhotos
	}

	if *request.EventEditSite.SportType != eventFromDB.SportType &&
		strings.Contains(*request.EventEditSite.URLPreview, "default") {
		defaultPhoto := a.getDefaultEventPhoto(*request.EventEditSite.SportType)
		request.EventEditSite.URLPreview = &defaultPhoto
		request.EventEditSite.URLPhotos = []string{defaultPhoto}
	}

	preResult := &models.FullEvent{
		ShortEvent: models.ShortEvent{
			ID:          request.EventID,
			CreatorID:   eventFromDB.CreatorID,
			SportType:   common.NewValWithFallback(request.EventEditSite.SportType, &eventFromDB.SportType),
			Address:     common.NewValWithFallback(request.EventEditSite.Address, &eventFromDB.Address),
			DateAndTime: common.NewValWithFallback(request.EventEditSite.DateAndTime, &eventFromDB.DateAndTime),
			Price:       request.EventEditSite.Price,
			GameLevels:  request.EventEditSite.GameLevels,
			Capacity:    request.EventEditSite.Capacity,
			URLPreview:  common.NewValWithFallback(request.EventEditSite.URLPreview, &eventFromDB.URLPreview),
			URLPhotos:   request.EventEditSite.URLPhotos,
		},
		Description:  request.EventEditSite.Description,
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

	err = a.onEventUpdate(ctx, preResult.ID)
	if err != nil {
		a.logger.WithCtx(ctx).Warnw("Unable to update bot event", "event_id", preResult.ID, "error", err)
	}

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

	err = a.onEventDelete(ctx, eventID)
	if err != nil {
		a.logger.WithCtx(ctx).Warnw("Unable to deleted bot event", "event_id", eventID, "error", err)
	}

	return nil
}

func (a *App) GetEvents(ctx context.Context) ([]models.ShortEvent, error) {
	return a.eventStorage.GetEvents(ctx)
}

func (a *App) FindEvents(ctx context.Context, filterParams *models.FilterParams) ([]models.ShortEvent, error) {
	if filterParams.Address != "" {
		func() {
			a.muFindByAddress.Lock()
			defer a.muFindByAddress.Unlock()

			latitude, longitude, err := a.getCoordinatesByAddress(ctx, filterParams.Address, UserAgentFind)
			if err != nil {
				a.logger.WithCtx(ctx).Errorf("to find address from=%s FindEvents: %v", filterParams.Address, err)
			} else {
				filterParams.AddressLatitude = &latitude
				filterParams.AddressLongitude = &longitude
			}

			time.Sleep(time.Millisecond * 1100)
		}()

	}

	return a.eventStorage.FindEvents(ctx, filterParams)
}

func (a *App) GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error) {
	return a.eventStorage.GetEvent(ctx, id)
}

func (a *App) SubscribeEvent(ctx context.Context, id uuid.UUID, userID *uuid.UUID, tgID *int64, subscribe bool) (*models.ResponseSubscribeEvent, error) {
	if userID == nil && tgID != nil {
		userFullFromTgID, err := a.authStorage.GetUserFullByTgID(ctx, *tgID)
		if err != nil {
			return nil, fmt.Errorf("to get user full by tg id: %w", err)
		}
		userID = &userFullFromTgID.ID
	}

	responseSubscribeEvent, err := a.eventStorage.SubscribeEvent(ctx, id, *userID, subscribe)
	if err != nil {
		return nil, fmt.Errorf("to subscribe event: %w", err)
	}

	err = a.onEventUpdate(ctx, responseSubscribeEvent.ID)
	if err != nil {
		a.logger.WithCtx(ctx).Warnw("Unable to update bot event", "event_id", responseSubscribeEvent.ID, "error", err)
	}

	return responseSubscribeEvent, nil
}

func (a *App) UserIsSubscribed(ctx context.Context, eventID uuid.UUID, reqParams *models.RequestUserIsSubscribedParams) (bool, error) {
	if reqParams.TgID != nil {
		userFullFromTgID, err := a.authStorage.GetUserFullByTgID(ctx, *reqParams.TgID)
		if err != nil {
			return false, fmt.Errorf("to get user full by tg id: %w", err)
		}
		reqParams.UserID = &userFullFromTgID.ID
	}

	eventFull, err := a.eventStorage.GetEvent(ctx, eventID)
	if err != nil {
		return false, fmt.Errorf("to get event: %w", err)
	}

	for _, subscriber := range eventFull.Subscribers {
		if subscriber == *reqParams.UserID {
			return true, nil
		}
	}

	return false, nil
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
