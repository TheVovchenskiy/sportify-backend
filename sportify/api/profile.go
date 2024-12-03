package api

import (
	"context"
	"errors"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/api"
	"github.com/go-pkgz/auth/token"
	"github.com/google/uuid"
	"net/http"
)

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

	if userInfo, err := token.GetUserInfo(r); err == nil {
		userIDFromToken, _ = uuid.Parse(userInfo.ID)
	}

	userFull, err := h.app.GetUserFullByUserID(ctx, profileUserID)
	if err != nil {
		h.handleGetProfile(ctx, w, err)
		return
	}

	result := models.MapUserFullToProfileAPI(userIDFromToken, userFull)

	models.WriteJSONResponse(w, result)
}
