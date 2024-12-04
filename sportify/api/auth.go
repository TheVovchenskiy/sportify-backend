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
	"time"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"

	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
)

func (h *Handler) handleCheck(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, db.ErrUserNotFound):
		models.WriteResponseError(w, models.NewResponseForbiddenErr("", db.ErrUserNotFound.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

func (h *Handler) WriteCheckResponse(ctx context.Context, w http.ResponseWriter, userInfo *token.User) {
	userFull, err := h.app.GetUserFullByUsername(ctx, userInfo.Name)
	if err != nil {
		err = fmt.Errorf("to get user full by name: %w", err)
		h.handleCheck(ctx, w, err)
		return
	}

	var responseCheck models.ResponseSuccessLogin

	responseCheck.Username = userInfo.Name
	responseCheck.UserID = userFull.ID

	models.WriteJSONResponse(w, responseCheck)
}

func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userInfo, err := token.GetUserInfo(r)
	if err != nil {
		err = fmt.Errorf("to get user info: %w", err)
		h.handleCheck(ctx, w, err)
		return
	}

	h.WriteCheckResponse(ctx, w, &userInfo)
}

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

	responseSuccessRegister, err := h.app.CreateUser(ctx, username, password)
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

func (h *Handler) handleTgAuth(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	// case errors.Is(errOutside, ErrRequestEventCreateSite):
	// 	models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	// case errors.Is(errOutside, app.ErrNotValidUsername):
	// 	models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	// case errors.Is(errOutside, app.ErrNotValidPassword):
	// 	models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	// case errors.Is(errOutside, app.ErrNotUniqueUsername):
	// 	models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

// func (h *Handler) TgAuth(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	reqBody, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		h.handleTgAuth(ctx, w, err)
// 		return
// 	}

// 	var tgRequestAuth models.TgRequestAuth
// 	err = json.Unmarshal(reqBody, &tgRequestAuth)
// 	if err != nil {
// 		err = fmt.Errorf("%w: %w", ErrRequestRegister, err)
// 		h.handleTgAuth(ctx, w, err)
// 		return
// 	}

// 	tgResponse, err := h.app.TgAuth(ctx, &tgRequestAuth)
// 	if err != nil {
// 		h.handleTgAuth(ctx, w, err)
// 		return
// 	}

// 	models.WriteJSONResponse(w, tgResponse)
// }

func (h *Handler) NewCredCheckFunc(ctx context.Context) provider.CredCheckerFunc {
	return h.app.NewCredCheckFunc(ctx)
}

func (h *Handler) handleLoginFromTg(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

var ErrRequestLoginFromTg = errors.New("не корректный запрос на логин из тг")

func (h *Handler) LoginUserFromTg(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var tgReqAuth models.TgRequestAuth

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleLoginFromTg(ctx, w, err)
		return
	}

	err = json.Unmarshal(reqBody, &tgReqAuth)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrRequestLoginFromTg, err)
		h.handleLoginFromTg(ctx, w, err)
		return
	}

	h.telegram.AddUpdate(&models.TgUpdateWrapper{
		TgUpdate:       tgReqAuth.TgUpdate,
		ExpirationTime: time.Now().Add(time.Second * 10),
	})

	err = h.app.LoginUserFromTg(ctx, &tgReqAuth)
	if err != nil {
		h.handleLoginFromTg(ctx, w, err)
	}

	err = h.telegram.GetResult(tgReqAuth.TgUpdate.Message.Chat.ID)
	if err != nil {
		h.handleLoginFromTg(ctx, w, err)
		return
	}

	models.WriteJSONResponse(w, "ok")
}
