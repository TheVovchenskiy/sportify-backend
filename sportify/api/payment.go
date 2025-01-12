package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/api"
)

func (h *Handler) handlePayEventError(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, ErrRequestEventPay):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, app.ErrPayFree):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", app.ErrPayFree.Error()))
	case errors.Is(errOutside, db.ErrNotFoundEvent):
		models.WriteResponseError(w, models.NewResponseNotFoundErr("", db.ErrNotFoundEvent.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

var ErrRequestEventPay = errors.New("Некорректный запрос на оплату события")

func (h *Handler) PayEvent(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleGetPaymentError(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, db.ErrNotFoundPayment):
		models.WriteResponseError(w, models.NewResponseNotFoundErr("", db.ErrNotFoundPayment.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.GetUUID(r, "id")
	if err != nil {
		h.handleGetPaymentError(ctx, w, err)
		return
	}

	payment, err := h.app.GetPayment(ctx, id)
	if err != nil {
		h.handleGetPaymentError(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, payment)
}
