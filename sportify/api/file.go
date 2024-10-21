package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/models"
)

const (
	MaxSizePhotoBytes = 50 * 1024 * 1024
)

func (h *Handler) handleUploadFileError(ctx context.Context, w http.ResponseWriter, errOutside error) {
	h.logger.WithCtx(ctx).Error(errOutside)

	switch {
	case errors.Is(errOutside, ErrParseFileBody):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	case errors.Is(errOutside, ErrToBigFile):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", ErrToBigFile.Error()))
	case errors.Is(errOutside, app.ErrWrongFormat):
		models.WriteResponseError(w, models.NewResponseBadRequestErr("", errOutside.Error()))
	default:
		models.WriteResponseError(w, models.NewResponseInternalServerErr("", models.InternalServerErrMessage))
	}
}

var (
	ErrParseFileBody = errors.New("ошибка парсинга файла")
	ErrToBigFile     = fmt.Errorf("слишком большой файл максимум %d", MaxSizePhotoBytes/1024/1024)
)

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleUploadFileError(ctx, w, err)
		return
	}

	lenBody := base64.StdEncoding.DecodedLen(len(body))
	if lenBody > MaxSizePhotoBytes {
		h.handleUploadFileError(ctx, w, ErrToBigFile)
		return
	}

	rawBody := make([]byte, lenBody)

	_, err = base64.StdEncoding.Decode(rawBody, body)
	if err != nil {
		h.handleUploadFileError(ctx, w, fmt.Errorf("%w: %s", ErrParseFileBody, err.Error()))
		return
	}

	url, err := h.app.SaveImage(ctx, rawBody)
	if err != nil {
		h.handleUploadFileError(ctx, w, fmt.Errorf("to save file: %w", err))
		return
	}

	models.WriteJSONResponse(w, models.ResponseUploadFile{URL: url})
}
