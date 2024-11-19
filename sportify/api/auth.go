package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/go-pkgz/auth/provider"
)

func (h *Handler) handleRegister(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, ErrRequestEventCreateSite):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, app.ErrNotValidUsername):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, app.ErrNotValidPassword):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, app.ErrNotUniqueUsername):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

var ErrRequestRegister = errors.New("не корректный запрос регистрации")

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestRegister models.RequestLogin

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleRegister(ctx, w, err)
		return
	}

	err = json.Unmarshal(reqBody, &requestRegister)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrRequestRegister, err)
		h.handleRegister(ctx, w, err)
		return
	}

	username, password, err := h.app.ValidateUsernameAndPassword(requestRegister.Username, requestRegister.PasswordRaw)
	if err != nil {
		h.handleRegister(ctx, w, err)
		return
	}

	responseSuccessRegister, err := h.app.CreateUser(ctx, username, password, nil)
	if err != nil {
		h.handleRegister(ctx, w, err)
		return
	}

	urlReqLogin, err := url.JoinPath("http://", h.url, h.apiPrefix, "/auth/my/login")
	if err != nil {
		err = fmt.Errorf("to join path: %w", err)
		h.handleRegister(ctx, w, err)
		return
	}

	h.logger.Info(urlReqLogin)

	requestLogin := models.RequestLogin{Username: username, PasswordRaw: password}

	requestLoginBody, err := json.Marshal(requestLogin)
	if err != nil {
		err = fmt.Errorf("to marshall requestLogin: %w", err)
		h.handleRegister(ctx, w, err)
		return
	}

	respLogin, err := http.Post(urlReqLogin, "application/json", bytes.NewReader(requestLoginBody))
	if err != nil {
		err = fmt.Errorf("to post login: %w", err)
		h.handleRegister(ctx, w, err)
		return
	}
	defer respLogin.Body.Close()

	if respLogin.StatusCode != http.StatusOK {
		h.handleRegister(ctx, w, fmt.Errorf("status code to login: %d", respLogin.StatusCode))
		return
	}

	for _, cookieVal := range respLogin.Header.Values("Set-Cookie") {
		w.Header().Add("Set-Cookie", cookieVal)
	}

	models.WriteJSONResponse(w, responseSuccessRegister)
}

func (h *Handler) NewCredCheckFunc(ctx context.Context) provider.CredCheckerFunc {
	return h.app.NewCredCheckFunc(ctx)
}
