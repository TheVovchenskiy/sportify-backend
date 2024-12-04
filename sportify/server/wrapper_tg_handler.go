package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/TheVovchenskiy/sportify-backend/db"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/api"
	"github.com/TheVovchenskiy/sportify-backend/pkg/middleware"

	"github.com/go-pkgz/auth/provider"
)

type StorageToken interface {
	GetUsername(ctx context.Context, token string) (string, error)
}

type WrapperTgHandler struct {
	storageToken StorageToken
	checkHandler middleware.CheckHandler
	*provider.TelegramHandler
}

func NewWrapperTgHandler(
	storageToken StorageToken,
	checkHandler middleware.CheckHandler,
	tgHandler *provider.TelegramHandler,
) *WrapperTgHandler {
	return &WrapperTgHandler{
		storageToken:    storageToken,
		checkHandler:    checkHandler,
		TelegramHandler: tgHandler,
	}
}

func (wth *WrapperTgHandler) LoginHandler(writer http.ResponseWriter, request *http.Request) {
	queryToken := request.URL.Query().Get("token")

	// in first case just get login token
	if queryToken == "" {
		wth.TelegramHandler.LoginHandler(writer, request)
		return
	}

	dummyWriter := httptest.NewRecorder()
	wth.TelegramHandler.LoginHandler(dummyWriter, request)

	if dummyWriter.Code != http.StatusOK {
		api.WriteFromDummyWriterToWriter(dummyWriter, writer)
		return
	}

	username, err := wth.storageToken.GetUsername(context.Background(), queryToken)
	if err != nil {
		fmt.Printf("WrapperTgHandler.LoginHandler get username: %+v", err)

		if errors.Is(err, db.ErrTokenNotFound) || errors.Is(err, db.ErrTokenExpired) {
			models.WriteResponseError(writer, models.NewResponseBadRequestErr(err.Error(), ""))
			return
		}

		models.WriteResponseError(writer, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
		return
	}

	api.WriteHeaderToWriter(dummyWriter.Header(), writer)
	wth.checkHandler.WriteCheckResponse(request.Context(), writer, request, username)
}
