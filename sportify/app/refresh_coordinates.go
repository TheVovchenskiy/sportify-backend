package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/models"
)

func requestURLOpenMap() string {
	return fmt.Sprintf("https://nominatim.openstreetmap.org/search" +
		"?&limit=1&accept-language=ru-RU&countrycodes=RU&format=jsonv2")
}

const userAgent = "SportifyApp/1.0"

func (a *App) getCoordinatesByAddress(ctx context.Context, address string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", requestURLOpenMap(), nil)
	if err != nil {
		return "", "", fmt.Errorf("to new request: %w", err)
	}

	values := req.URL.Query()
	values.Add("q", address)
	req.URL.RawQuery = values.Encode()

	a.logger.Infof(req.URL.String())

	req.Header.Set("User-Agent", userAgent)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("to do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("to read body: %w", err)
	}

	a.logger.Info("GET COORDINATES: ", address, string(body))

	var coordinates []models.ResponseOpenMapCoordinates

	err = json.Unmarshal(body, &coordinates)
	if err != nil {
		return "", "", fmt.Errorf("to unmarshal body: %w", err)
	}

	if len(coordinates) != 1 {
		return "", "", fmt.Errorf("invalid coordinates count: %d", len(coordinates))
	}

	return coordinates[0].Latitude, coordinates[0].Longitude, nil
}

func (a *App) RefreshCoordinates(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(period)

			params, err := models.ParseFilterParams(make(url.Values))
			if err != nil {
				a.logger.WithCtx(ctx).Error(err)
				continue
			}

			events, err := a.eventStorage.FindEvents(ctx, params)
			if err != nil {
				a.logger.WithCtx(ctx).Error(err)
				continue
			}

			for _, event := range events {
				if event.Latitude == nil || event.Longitude == nil {
					latitude, longitude, err := a.getCoordinatesByAddress(ctx, event.Address)
					if err != nil {
						a.logger.WithCtx(ctx).Error(err)
						continue
					}

					err = a.eventStorage.SetCoordinates(ctx, latitude, longitude, event.ID)
					if err != nil {
						a.logger.WithCtx(ctx).Error(err)
						continue
					}

					a.logger.Infof("set coordinates for event %s: %s, %s", event.ID.String(), latitude, longitude)
				}

				// limit of open map
				time.Sleep(time.Second * 2)
			}
		}
	}
}
