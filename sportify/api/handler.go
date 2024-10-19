package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
	"github.com/TheVovchenskiy/sportify-backend/pkg/mylogger"

	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type App interface {
	AddEvent(ctx context.Context, event models.FullEvent) error
	GetEvents(ctx context.Context) ([]models.ShortEvent, error)
	GetEvent(ctx context.Context, id uuid.UUID) (*models.FullEvent, error)
	SubscribeEvent(ctx context.Context, d uuid.UUID, userID uuid.UUID, subscribe bool) (*models.ResponseSubscribeEvent, error)
	DetectEventMessage(text string, regexps []string, minMatches int) (bool, error)
}

var _ App = (*app.App)(nil)

type Handler struct {
	logger *mylogger.MyLogger
	app    App
}

func NewHandler(app App, logger *mylogger.MyLogger) Handler {
	return Handler{app: app, logger: logger}
}

func (h *Handler) handleGetEventsError(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)
	models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	events, err := h.app.GetEvents(ctx)
	if err != nil {
		h.handleGetEventsError(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, events)
}

var ErrInvalidEventID = errors.New("не верный event id")

func (h *Handler) handleGetEventError(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, ErrInvalidEventID):
		models.WriteResponseError(w, models.NewResponseBadRequest("", ErrInvalidEventID.Error()))
	case errors.Is(errOutside, db.ErrNotFoundEvent):
		models.WriteResponseError(w, models.NewResponseBadRequest("", db.ErrNotFoundEvent.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	preEventID := chi.URLParam(r, "id")

	eventID, err := uuid.Parse(preEventID)
	if err != nil {
		err = fmt.Errorf("eventID %s: %w", err.Error(), ErrInvalidEventID)
		h.handleGetEventError(ctx, w, err)
		return
	}

	event, err := h.app.GetEvent(ctx, eventID)
	if err != nil {
		h.handleGetEventError(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, event)
}

var ErrRequestSubscribeEvent = errors.New("не корректный запрос subscribe event")

func (h *Handler) handleSubscribeEventErr(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, ErrInvalidEventID):
		models.WriteResponseError(w, models.NewResponseBadRequest("", ErrInvalidEventID.Error()))
	case errors.Is(errOutside, ErrRequestSubscribeEvent):
		models.WriteResponseError(w, models.NewResponseBadRequest("", ErrRequestSubscribeEvent.Error()))
	case errors.Is(errOutside, db.ErrNotFoundEvent):
		models.WriteResponseError(w, models.NewResponseBadRequest("", db.ErrNotFoundEvent.Error()))
	case errors.Is(errOutside, models.ErrAllBusy):
		models.WriteResponseError(w, models.NewResponseBadRequest("", models.ErrAllBusy.Error()))
	case errors.Is(errOutside, models.ErrFoundSubscriber):
		models.WriteResponseError(w, models.NewResponseBadRequest("", models.ErrFoundSubscriber.Error()))
	case errors.Is(errOutside, models.ErrNotFoundSubscriber):
		models.WriteResponseError(w, models.NewResponseBadRequest("", models.ErrNotFoundSubscriber.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

func (h *Handler) SubscribeEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	preEventID := chi.URLParam(r, "id")

	eventID, err := uuid.Parse(preEventID)
	if err != nil {
		err = fmt.Errorf("eventID %s: %w", err.Error(), ErrInvalidEventID)

		h.handleSubscribeEventErr(ctx, w, err)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleSubscribeEventErr(ctx, w, err)
		return
	}

	var reqSubEvent models.RequestSubscribeEvent
	err = json.Unmarshal(reqBody, &reqSubEvent)
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), ErrRequestSubscribeEvent)

		h.handleSubscribeEventErr(ctx, w, err)
		return
	}

	responseSubscribeEvent, err := h.app.SubscribeEvent(ctx, eventID, reqSubEvent.UserID, reqSubEvent.SubscribeFlag)
	if err != nil {
		h.handleSubscribeEventErr(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, responseSubscribeEvent)
}

func templateBodyYaGPT(folderID, text string) string {
	return fmt.Sprintf("{\n  \"modelUri\": \"gpt://%s/yandexgpt-lite\",\n  \"completionOptions\": {\n    \"stream\": false,\n    \"temperature\": 0.1,\n    \"maxTokens\": \"1000\"\n  },\n  \"messages\": [\n    {\n      \"role\": \"system\",\n      \"text\": \"Тебе нужно распарсить из сообщения информацию в формате json:\\n{\\\"cost\\\": \\\"200\\\",\\n\\\"date\\\": \\\"12.10\\\",\\n\\\"start_time\\\": \\\"18:00\\\",\\n\\\"end_time\\\": \\\"18:00\\\",\\n\\\"location\\\": \\\"г. Москва ул. 50-Летия Победы д.22 или м. Белорусская\\\"}\\n\\nСтрого соблюдай требования: поле \\\"cost\\\" должно быть числом - количеством рублей,\\nполе \\\"date\\\" 20.10 именно в формате месяц.день год указывать не нужно!,\\nполе \\\"start_time\\\" именно часы:минуты,\\nполе \\\"end_time\\\" 18:00 именно часы:минуты,\\nполе \\\"location\\\" любую информацию про местоположение.\\n\\nЕсли какое-то поле не получилось найти, оставь поле пустым вот так \\\"\\\".\\n\"\n    },\n    {\n      \"role\": \"user\",\n      \"text\": \"%s\"\n    }\n  ]\n}", folderID, text)
}

type responseYaGPT struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

func (r *responseYaGPT) getText() string {
	return strings.Trim(r.Result.Alternatives[0].Message.Text, "`\n")
}

type eventYaGPT struct {
	Cost      int    `json:"cost"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
}

func eventFromYaGPT(text []byte) (*models.FullEvent, error) {
	eventYa := eventYaGPT{}

	err := json.Unmarshal(text, &eventYa)
	if err != nil {
		return nil, err
	}

	var result models.FullEvent

	idxDot := strings.Index(eventYa.Date, ".")
	eventYa.Date = eventYa.Date[idxDot+1:] + "." + eventYa.Date[:idxDot]

	date, err := time.Parse("01.02", eventYa.Date)
	if err != nil {
		return nil, err
	}

	result.Date = time.Date(2024, date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	startTime, err := time.Parse("15:04", eventYa.StartTime)
	if err != nil {
		return nil, err
	}

	result.StartTime = time.Date(2024, date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)

	if eventYa.EndTime == "" {
		endTime, err := time.Parse("15:04", eventYa.EndTime)
		if err != nil {
			return nil, err
		}

		result.EndTime = common.Ref(time.Date(2024, date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC))
	}

	result.Address = eventYa.Location
	result.ID = uuid.New()
	result.Price = common.Ref(eventYa.Cost)
	result.SportType = models.SportTypeFootball
	result.URLPreview = "http://127.0.0.1:8080/img/default_football.jpeg"

	return &result, nil
}

func (h *Handler) requestToYaGPT(text string) (*models.FullEvent, error) {
	folderID := os.Getenv("FOLDER_ID")
	iamToken := os.Getenv("IAM_TOKEN")

	body := templateBodyYaGPT(folderID, text)

	req, err := http.NewRequest(http.MethodPost, "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header["x-folder-id"] = []string{folderID}
	req.Header["Authorization"] = []string{"Bearer " + iamToken}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responseYA responseYaGPT

	err = json.Unmarshal(resBody, &responseYA)
	if err != nil {
		return nil, err
	}

	return eventFromYaGPT([]byte(responseYA.getText()))
}

func (h *Handler) handleTryCreateEventErr(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, db.ErrEventAlreadyExist):
		models.WriteResponseError(w, models.NewResponseBadRequest("", db.ErrEventAlreadyExist.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

func (h *Handler) TryCreateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tgMessage models.TgMessage

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleTryCreateEventErr(ctx, w, err)
		return
	}

	err = json.Unmarshal(reqBody, &tgMessage)
	if err != nil {
		h.handleTryCreateEventErr(ctx, w, err)
		return
	}

	if ok, err := h.app.DetectEventMessage(tgMessage.RawMessage, app.SportEventRegExps, 3); !ok || err != nil {
		fmt.Println("err detect: ", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	fullEvent, err := h.requestToYaGPT(strings.ReplaceAll(tgMessage.RawMessage, "\n", " "))
	if err != nil {
		h.handleGetEventError(ctx, w, err)
		return
	}

	err = h.app.AddEvent(ctx, *fullEvent)
	if err != nil {
		h.handleGetEventError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
