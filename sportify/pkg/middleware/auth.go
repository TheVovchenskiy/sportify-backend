package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/api"

	"github.com/go-pkgz/auth/token"
)

func PostOnlyRestriction(url string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.URL.Path, url) && request.Method != http.MethodPost {
			models.WriteResponseError(writer, models.ResponseErr{
				StatusCode: http.StatusMethodNotAllowed,
				ErrMessage: "разрешен только метод POST",
			})
		}

		next.ServeHTTP(writer, request)
	})
}

type CheckHandler interface {
	WriteCheckResponse(ctx context.Context, w http.ResponseWriter, r *http.Request, username string)
}

func ConvertLoginResponseToCheck(checkHandler CheckHandler, prev http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.URL.Path, "my/login") {
			dummyWriter := httptest.NewRecorder()

			prev.ServeHTTP(dummyWriter, request)

			if dummyWriter.Code != http.StatusOK {
				api.WriteFromDummyWriterToWriter(dummyWriter, writer)
				return
			}

			var userInfo token.User

			err := json.Unmarshal(dummyWriter.Body.Bytes(), &userInfo)
			if err != nil {
				fmt.Printf("ConvertLoginResponseToCheck to umarshall: %s\n", err.Error())
				models.WriteResponseError(writer, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
				return
			}

			api.WriteHeaderToWriter(dummyWriter.Header(), writer)
			checkHandler.WriteCheckResponse(request.Context(), writer, request, userInfo.Name)
		} else {
			prev.ServeHTTP(writer, request)
		}
	})
}

var mapErrReplace = map[string]string{
	"Unauthorized\n":             "Вы не авторизованы",
	"incorrect user or password": "Не верный логин или пароль",

	// from telegram check token
	"request is not found":        "Произошла какая-то ошибка авторизации, попробуйте снова",
	"request expired":             "Произошла какая-то ошибка авторизации, попробуйте снова",
	"request is not verified yet": "Подтвердите авторизацию в телеграм боте",
}

func ConvertErrUnknownToOurType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dummyWriter := httptest.NewRecorder()
		next.ServeHTTP(dummyWriter, r)

		for key, values := range dummyWriter.Header() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		body := string(dummyWriter.Body.Bytes())

		if dummyWriter.Code >= 400 {
			for oldErr, newErr := range mapErrReplace {
				body = strings.ReplaceAll(body, oldErr, newErr)
			}

			if !strings.Contains(body, `"error":`) {
				models.WriteResponseError(w, models.ResponseErr{StatusCode: dummyWriter.Code, ErrMessage: body})
				return
			}
		}

		w.WriteHeader(dummyWriter.Code)
		w.Write([]byte(body))
	})
}
