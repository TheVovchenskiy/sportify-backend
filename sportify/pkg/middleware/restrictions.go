package middleware

import (
	"github.com/TheVovchenskiy/sportify-backend/models"
	"net/http"
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
