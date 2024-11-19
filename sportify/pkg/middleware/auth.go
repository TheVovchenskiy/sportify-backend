package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/go-pkgz/auth/token"
	"net/http"
	"net/http/httptest"
	"strings"
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
	WriteCheckResponse(ctx context.Context, w http.ResponseWriter, userInfo *token.User)
}

func ConvertLoginResponseToCheck(handler CheckHandler, prev http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.URL.Path, "my/login") {
			ctx := request.Context()
			dummyWriter := httptest.NewRecorder()

			prev.ServeHTTP(dummyWriter, request)

			if dummyWriter.Code != http.StatusOK {
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(dummyWriter.Code)
				for key, values := range dummyWriter.Header() {
					for _, value := range values {
						writer.Header().Add(key, value)
					}
				}
				writer.Write(dummyWriter.Body.Bytes())
				return
			}

			var userInfo token.User

			err := json.Unmarshal(dummyWriter.Body.Bytes(), &userInfo)
			if err != nil {
				fmt.Printf("ConvertLoginResponseToCheck to umarshall: %s\n", err.Error())
				models.WriteResponseError(writer, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
				return
			}

			handler.WriteCheckResponse(ctx, writer, &userInfo)
		} else {
			prev.ServeHTTP(writer, request)
		}
	})
}
