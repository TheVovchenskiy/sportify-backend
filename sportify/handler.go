package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
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
	w.WriteHeader(http.StatusOK)
}
