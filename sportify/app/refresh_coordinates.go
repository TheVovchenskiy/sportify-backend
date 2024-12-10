package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TheVovchenskiy/sportify-backend/pkg/reformat_url_open_map"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"sync"
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

	a.logger.Debug(req.URL.String())

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

	a.logger.Debug("GET COORDINATES: ", address, string(body))

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
	type coordinates struct {
		ID      uuid.UUID
		Address string
	}

	var (
		queueCoordinates       []coordinates
		isIDInQueueCoordinates = make(map[uuid.UUID]struct{})
		muQueue                sync.RWMutex
	)

	go func() {
		tickerQueueCoordinates := time.NewTicker(time.Millisecond * 1500)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerQueueCoordinates.C:
				if len(queueCoordinates) != 0 {
					func() {
						muQueue.Lock()
						defer muQueue.Unlock()

						curCoordinate := queueCoordinates[0]
						queueCoordinates = queueCoordinates[1:]
						delete(isIDInQueueCoordinates, curCoordinate.ID)

						latitude, longitude, err := a.getCoordinatesByAddress(ctx, curCoordinate.Address)
						if err != nil {
							a.logger.WithCtx(ctx).Error(
								fmt.Sprintf("address: %s, err: %s", curCoordinate.Address, err.Error()),
							)
							queueCoordinates = append(queueCoordinates, curCoordinate)
							isIDInQueueCoordinates[curCoordinate.ID] = struct{}{}
							return
						}

						err = a.eventStorage.SetCoordinates(ctx, latitude, longitude, curCoordinate.ID)
						if err != nil {
							a.logger.WithCtx(ctx).Error(err)
							queueCoordinates = append(queueCoordinates, curCoordinate)
							isIDInQueueCoordinates[curCoordinate.ID] = struct{}{}
							return
						}

						a.logger.Infof("set coordinates for event %s: %s, %s", curCoordinate.ID.String(), latitude, longitude)
					}()
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
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

			// TODO optimize from read all db to read only WHERE coordinates IS NULL
			events, err := a.eventStorage.FindEvents(ctx, params)
			if err != nil {
				a.logger.WithCtx(ctx).Error(err)
				continue
			}

			for _, event := range events {
				muQueue.Lock()
				_, ok := isIDInQueueCoordinates[event.ID]
				if event.Latitude == nil && event.Longitude == nil && !ok {
					address := reformat_url_open_map.ReformatURLOpenMap(event.Address)
					queueCoordinates = append(queueCoordinates, coordinates{
						ID:      event.ID,
						Address: address,
					})
					isIDInQueueCoordinates[event.ID] = struct{}{}
				}
				muQueue.Unlock()
			}
		}
	}
}
