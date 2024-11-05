package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/models"
)

func (h *Handler) handlePayEventError(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, ErrRequestEventPay):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, app.ErrPayFree):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", app.ErrPayFree.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

var ErrRequestEventPay = errors.New("некорректный запрос на оплату события")

func (h *Handler) PayEvent(r *http.Request, w http.ResponseWriter) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.handlePayEventError(ctx, w, err)
		return
	}

	var requestEventPay models.RequestEventPay

	err = json.Unmarshal(body, &requestEventPay)
	if err != nil {
		errOutside := fmt.Errorf("%w: %s", ErrRequestEventPay, err.Error())

		h.handlePayEventError(ctx, w, errOutside)
		return
	}

	responseEventPay, err := h.app.PayEvent(ctx, &requestEventPay)
	if err != nil {
		h.handlePayEventError(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, responseEventPay)
}
