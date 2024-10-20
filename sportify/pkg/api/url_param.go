package api

import (
	"errors"
	"fmt"
	"net/http"

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
