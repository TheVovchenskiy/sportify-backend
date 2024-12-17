package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	chi "github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var ErrInvalidUUID = errors.New("не верный uuid")

func GetUUID(r *http.Request, param string) (uuid.UUID, error) {
	preUUID := chi.URLParam(r, param)

	result, err := uuid.Parse(preUUID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: параметр %s: %s", ErrInvalidUUID, param, err.Error())
	}

	return result, nil
}

func GetInt64(r *http.Request, param string) (int64, error) {
	preInt64 := chi.URLParam(r, param)

	result, err := strconv.ParseInt(preInt64, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: параметр %s: %s", ErrInvalidUUID, param, err.Error())
	}

	return result, nil
}
