package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Handler struct {
	Storage EventStorage
}

func (h *Handler) HandleError(w http.ResponseWriter, err error) {
	log.Println(err)

	switch {
	case errors.Is(err, ErrEventAlreadyExist):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrEventAlreadyExist.Error()))
	case errors.Is(err, ErrNotFoundEvent):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(ErrNotFoundEvent.Error()))
	case errors.Is(err, ErrNotFoundSubscriber):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(ErrNotFoundSubscriber.Error()))
	case errors.Is(err, ErrAllBusy):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrAllBusy.Error()))
	case errors.Is(err, ErrInvalidUUID):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	case errors.Is(err, ErrRequestSubscribeEvent):
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error in server"))
	}
}

var ErrInvalidUUID = errors.New("invalid uuid")

func (h *Handler) GetEvents(w http.ResponseWriter, _ *http.Request) {
	events, err := h.Storage.GetEvents()
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	body, err := json.Marshal(events)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.Write(body)
}

func (h *Handler) GetEvent(w http.ResponseWriter, r *http.Request) {
	preEventID := chi.URLParam(r, "id")

	eventID, err := uuid.Parse(preEventID)
	if err != nil {
		err = fmt.Errorf("eventID %s: %w", err.Error(), ErrInvalidUUID)

		h.HandleError(w, err)
		return
	}

	event, err := h.Storage.GetEvent(eventID)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	body, err := json.Marshal(event)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.Write(body)
}

var ErrRequestSubscribeEvent = errors.New("invalid request subscribe event")

func (h *Handler) SubscribeEvent(w http.ResponseWriter, r *http.Request) {
	preEventID := chi.URLParam(r, "id")

	eventID, err := uuid.Parse(preEventID)
	if err != nil {
		err = fmt.Errorf("eventID %s: %w", err.Error(), ErrInvalidUUID)

		h.HandleError(w, err)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	var reqSubEvent RequestSubscribeEvent
	err = json.Unmarshal(reqBody, &reqSubEvent)
	if err != nil {
		err = fmt.Errorf("%s: %w", err, ErrRequestSubscribeEvent)

		h.HandleError(w, err)
		return
	}

	responseSubscribeEvent, err := h.Storage.SubscribeEvent(eventID, reqSubEvent.UserID, reqSubEvent.SubscribeFlag)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	respBody, err := json.Marshal(responseSubscribeEvent)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.Write(respBody)
}

func templateBodyYaGPT(folderID, text string) string {
	return fmt.Sprintf("{\n  \"modelUri\": \"gpt://%s/yandexgpt-lite\",\n  \"completionOptions\": {\n    \"stream\": false,\n    \"temperature\": 0.1,\n    \"maxTokens\": \"1000\"\n  },\n  \"messages\": [\n    {\n      \"role\": \"system\",\n      \"text\": \"Тебе нужно распарсить из сообщения информацию в формате json:\\n{\\\"cost\\\": \\\"200\\\",\\n\\\"date\\\": \\\"12.10\\\",\\n\\\"start_time\\\": \\\"18:00\\\",\\n\\\"end_time\\\": \\\"18:00\\\",\\n\\\"location\\\": \\\"г. Москва ул. 50-Летия Победы д.22 или м. Белорусская\\\"}\\n\\nСтрого соблюдай требования: поле \\\"cost\\\" должно быть числом - количеством рублей,\\nполе \\\"date\\\" 20.10 именно в формате месяц.день год указывать не нужно!,\\nполе \\\"start_time\\\" именно часы:минуты,\\nполе \\\"end_time\\\" 18:00 именно часы:минуты,\\nполе \\\"location\\\" любую информацию про местоположение.\\n\\nЕсли какое-то поле не получилось найти, оставь поле пустым вот так \\\"\\\".\\n\"\n    },\n    {\n      \"role\": \"user\",\n      \"text\": \"%s\"\n    }\n  ]\n}", folderID, text)
}

// {
// "result": {
// "alternatives": [
// {
// "message": {
// "role": "assistant",
// "text": "Ламинат подходит для укладки на кухне и в детской комнате. Он не боится влажности и механических повреждений, благодаря защитному слою, состоящему из меланиновых плёнок толщиной 0.2 мм, и обработанным воском замкам."
// },
// "status": "ALTERNATIVE_STATUS_TRUNCATED_FINAL"
// }
// ],
// "usage": {
// "inputTextTokens": "67",
// "completionTokens": "50",
// "totalTokens": "117"
// },
// "modelVersion": "06.12.2023"
// }
// }
type ResponseYaGPT struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

func (r *ResponseYaGPT) GetText() string {
	return strings.Trim(r.Result.Alternatives[0].Message.Text, "`\n")
}

//"cost": 200,
//"date": "20.10",
//"start_time": "18:00",
//"end_time": "18:00",
//"location": "г. Москва ул. 50-Летия Победы д.22 или м. Белорусская"
//

type EventYaGPT struct {
	Cost      int    `json:"cost"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Location  string `json:"location"`
}

func EventFromYaGPT(text []byte) (*FullEvent, error) {
	eventYa := EventYaGPT{}

	err := json.Unmarshal(text, &eventYa)
	if err != nil {
		return nil, err
	}

	var result FullEvent

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

	result.StartTime =
		time.Date(2024, date.Month(), date.Day(), startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)

	if eventYa.EndTime == "" {
		endTime, err := time.Parse("15:04", eventYa.EndTime)
		if err != nil {
			return nil, err
		}

		result.EndTime =
			Ref(time.Date(2024, date.Month(), date.Day(), endTime.Hour(), endTime.Minute(), 0, 0, time.UTC))
	}

	result.Address = eventYa.Location
	result.ID = uuid.New()
	result.Price = Ref(eventYa.Cost)
	result.SportType = TypeFootball
	result.PreviewURL = "http://127.0.0.1:8080/img/default_football.jpeg"

	return &result, nil
}

func (h *Handler) RequestToYaGPT(text string) (*FullEvent, error) {
	folderID := os.Getenv("FOLDER_ID")
	iamToken := os.Getenv("IAM_TOKEN")

	//fmt.Println(folderID, iamToken)

	body := templateBodyYaGPT(folderID, text)

	fmt.Println(body)

	req, err := http.NewRequest(http.MethodPost, "https://llm.api.cloud.yandex.net/foundationModels/v1/completion", bytes.NewReader([]byte(body)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header["x-folder-id"] = []string{folderID}
	req.Header["Authorization"] = []string{"Bearer " + iamToken}

	log.Printf("request %+v\n", req)

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

	log.Printf("response %+v\n", res)

	log.Println()
	log.Println()
	log.Println("body", string(resBody))

	var responseYA ResponseYaGPT

	err = json.Unmarshal(resBody, &responseYA)
	if err != nil {
		return nil, err
	}

	return EventFromYaGPT([]byte(responseYA.GetText()))
}

func (h *Handler) TryCreateEvent(w http.ResponseWriter, r *http.Request) {
	var tgMessage TgMessage

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	err = json.Unmarshal(reqBody, &tgMessage)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	log.Printf("%+v", tgMessage)

	// TODO detect

	if ok, err := detect(tgMessage.RawMessage, detectRegExps, 3); !ok || err != nil {
		fmt.Println("err detect: ", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	fullEvent, err := h.RequestToYaGPT(strings.ReplaceAll(tgMessage.RawMessage, "\n", " "))
	if err != nil {
		h.HandleError(w, err)
		return
	}

	err = h.Storage.AddEvent(*fullEvent)
	if err != nil {
		h.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
