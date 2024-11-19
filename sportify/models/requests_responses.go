package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	"github.com/google/uuid"
)

type ResponseOpenMapCoordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

type ResponsesPayment struct {
	PaymentStatus PaymentStatus `json:"payment_status"`
}

type RequestEventPay struct {
	RedirectURL string    `json:"redirect_url"`
	UserID      uuid.UUID `json:"user_id"`
	EventID     uuid.UUID `json:"event_id"`
}

type ResponseEventPay struct {
	ID              uuid.UUID `json:"id"`
	ConfirmationURL string    `json:"confirmation_url"`
}

type RequestEventDelete struct {
	UserID uuid.UUID `json:"user_id"`
}

type RequestEventEditSite struct {
	EventID       uuid.UUID     `json:"-"`
	UserID        uuid.UUID     `json:"user_id"`
	EventEditSite EventEditSite `json:"event_edit"`
}

type EventEditSite struct {
	SportType   *SportType   `json:"sport_type"`
	Address     *string      `json:"address"`
	DateAndTime *DateAndTime `json:"date_time"`
	Price       *int         `json:"price"`
	GameLevels  []GameLevel  `json:"game_level"`
	Description *string      `json:"description"`
	Capacity    *int         `json:"capacity"`
	URLPreview  *string      `json:"preview"`
	URLPhotos   []string     `json:"photos"`
}

type EventCreateSite struct {
	SportType   SportType   `json:"sport_type"`
	Address     string      `json:"address"`
	DateAndTime DateAndTime `json:"date_time"`
	Price       *int        `json:"price"`
	GameLevels  []GameLevel `json:"game_level"`
	Description *string     `json:"description"`
	Capacity    *int        `json:"capacity"`
	URLPreview  string      `json:"preview"`
	URLPhotos   []string    `json:"photos"`
}

type TgParams struct {
	UserID *int64  `json:"user_id"`
	ChatID *string `json:"chat_id"`
}

type RequestEventCreateSite struct {
	UserID      uuid.UUID       `json:"user_id"`
	Tg          *TgParams       `json:"tg,omitempty"`
	CreateEvent EventCreateSite `json:"event_create"`
}

type RequestSubscribeEvent struct {
	SubscribeFlag bool      `json:"sub"`
	UserID        uuid.UUID `json:"user_id"`
}

func WriteJSONResponse(w http.ResponseWriter, response any) {
	logger, err := mylogger.Get()
	if err != nil {
		fmt.Println(err)
		logger = mylogger.NewNop()
	}

	body, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(InternalServerErrMessage))
		if err != nil {
			logger.Error("error write resp InternalServerErrMessage")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		logger.Error("error write resp body")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type ResponseEventDelete struct {
	Status string `json:"status"`
}

func NewResponseEventDelete() ResponseEventDelete {
	return ResponseEventDelete{Status: "ok"}
}

type ResponseUploadFile struct {
	URL string `json:"url"`
}

type ResponseSubscribeEvent struct {
	ID          uuid.UUID   `json:"id"`
	Capacity    *int        `json:"capacity"`
	Busy        int         `json:"busy"`
	Subscribers []uuid.UUID `json:"subscribers_id"`
}

var (
	ErrAllBusy            = errors.New("все места заняты")
	ErrFoundSubscriber    = errors.New("вы уже подписаны на это событие")
	ErrNotFoundSubscriber = errors.New("не найден подписчик события")
)

func (r *ResponseSubscribeEvent) AddSubscriber(id uuid.UUID) error {
	if r.Capacity != nil && *r.Capacity <= r.Busy {
		return ErrAllBusy
	}

	_, isFound := common.Find(r.Subscribers, func(item uuid.UUID) bool {
		return item == id
	})
	if isFound {
		return ErrFoundSubscriber
	}

	r.Subscribers = append(r.Subscribers, id)
	r.Busy = len(r.Subscribers)

	return nil
}

func (r *ResponseSubscribeEvent) RemoveSubscriber(id uuid.UUID) error {
	for i, v := range r.Subscribers {
		if v == id {
			r.Subscribers = append(r.Subscribers[:i], r.Subscribers[i+1:]...)

			r.Busy = len(r.Subscribers)

			return nil
		}
	}

	return ErrNotFoundSubscriber
}

type ResponseErr struct {
	StatusCode int    `json:"-"`
	ErrName    string `json:"error_name"`
	ErrMessage string `json:"error"`
}

func NewResponseBadRequestErr(name, message string) ResponseErr {
	return ResponseErr{
		StatusCode: http.StatusBadRequest,
		ErrName:    name,
		ErrMessage: message,
	}
}

func NewResponseForbiddenErr(name, message string) ResponseErr {
	return ResponseErr{
		StatusCode: http.StatusForbidden,
		ErrName:    name,
		ErrMessage: message,
	}
}

func NewResponseNotFoundErr(name, message string) ResponseErr {
	return ResponseErr{
		StatusCode: http.StatusNotFound,
		ErrName:    name,
		ErrMessage: message,
	}
}

func NewResponseInternalServerErr(name, message string) ResponseErr {
	return ResponseErr{
		StatusCode: http.StatusInternalServerError,
		ErrName:    name,
		ErrMessage: message,
	}
}

const InternalServerErrMessage = "Внутренняя ошибка на сервере"

func WriteResponseError(w http.ResponseWriter, responseError ResponseErr) {
	body, err := json.Marshal(responseError)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalServerErrMessage))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseError.StatusCode)
	w.Write(body)
}
