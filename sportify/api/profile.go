package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/api"
	"github.com/go-pkgz/auth/token"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) AddTokenServiceProvider(tokenService *token.Service) {
	h.tokenService = tokenService
}

func (h *Handler) handleGetProfile(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, api.ErrInvalidUUID):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, db.ErrUserNotFound):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", db.ErrUserNotFound.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	profileUserID, err := api.GetUUID(r, "id")
	if err != nil {
		h.handleGetProfile(ctx, w, err)
		return
	}

	var userIDFromToken uuid.UUID

	claims, _, err := h.tokenService.Get(r)
	if err == nil && claims.User != nil {
		userID := strings.TrimPrefix(strings.TrimPrefix(claims.User.ID, "telegram_"), "my_")
		userIDFromToken, err = uuid.Parse(userID)
		if err != nil {
			h.logger.Errorf("to parse uuid from claims=%s: %s", userID, err.Error())
		}
	}

	userFull, err := h.app.GetUserFullByUserID(ctx, profileUserID)
	if err != nil {
		h.handleGetProfile(ctx, w, err)
		return
	}

	result := models.MapUserFullToProfileAPI(userIDFromToken, userFull)

	models.WriteJSONResponse(w, result)
}

func (h *Handler) handleUpdateProfile(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, api.ErrInvalidUUID):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, ErrRequestUpdateProfile):
		models.WriteResponseError(w, models.NewResponseBadRequestErr(errOutside.Error(), ErrRequestUpdateProfile.Error()))
	case errors.Is(errOutside, app.ErrValidationRequestUpdateProfile):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

var ErrRequestUpdateProfile = errors.New("не корректный запрос на обновление")

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	profileUserID, err := api.GetUUID(r, "user_id")
	if err != nil {
		h.handleGetProfile(ctx, w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleGetProfile(ctx, w, err)
		return
	}

	var reqUpdateProfile models.RequestUpdateProfile

	err = json.Unmarshal(body, &reqUpdateProfile)
	if err != nil {
		errOutside := fmt.Errorf("%w: %s", ErrRequestUpdateProfile, err.Error())

		h.handleGetProfile(ctx, w, errOutside)
		return
	}

	err = h.app.UpdateProfile(ctx, profileUserID, reqUpdateProfile)
	if err != nil {
		h.handleGetProfile(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, models.NewResponseUpdateProfile("https://"+h.domain, profileUserID))
}
