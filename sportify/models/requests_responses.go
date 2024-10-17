package models

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type RequestSubscribeEvent struct {
	SubscribeFlag bool      `json:"sub"`
	UserID        uuid.UUID `json:"user_id"`
}

func WriteJSONResponse(w http.ResponseWriter, response any) {
	body, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalServerErrMessage))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}

type ResponseSubscribeEvent struct {
	ID          uuid.UUID   `json:"id"`
	Capacity    *int        `json:"capacity"`
	Busy        int         `json:"busy"`
	Subscribers []uuid.UUID `json:"subscribers_id"`
}

type ResponseErr struct {
	StatusCode int    `json:"-"`
	ErrName    string `json:"error_name"`
	ErrMessage string `json:"error_message"`
}

func NewResponseBadRequest(name, message string) ResponseErr {
	return ResponseErr{
		StatusCode: http.StatusBadRequest,
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

	w.WriteHeader(responseError.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
