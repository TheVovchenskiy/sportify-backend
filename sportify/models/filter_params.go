package models

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/TheVovchenskiy/sportify-backend/pkg/common"

	"github.com/google/uuid"
)

type FilterParams struct {
	// Block public usage (from query params)

	SportTypes []SportType
	GameLevels []GameLevel
	DateStarts []string
	PriceMin   *int
	PriceMax   *int
	FreePlaces *int
	Address    string
	OrderBy    string
	SortOrder  string

	// Block inner usage

	CreatorID     *uuid.UUID
	SubscriberIDs []uuid.UUID

	// DateExpression is representation of WHERE statement
	// you can use squirrel.Eq and another with similar sense
	DateExpression any

	AddressLatitude, AddressLongitude *string
}

//nolint:cyclop
func ParseFilterParams(query url.Values) (*FilterParams, error) {
	params := &FilterParams{ //nolint:exhaustruct
		SportTypes: common.Map[string, SportType](
			func(s string) SportType {
				return SportType(s)
			},
			query["sport_type"]),
		GameLevels: common.Map[string, GameLevel](
			func(s string) GameLevel {
				return GameLevel(s)
			}, query["game_level"]),
		DateStarts: query["date_start"],
		Address:    query.Get("address"),
		OrderBy:    query.Get("order_by"),
		SortOrder:  query.Get("sort_order"),
	}

	if priceMinStr := query.Get("price_min"); priceMinStr != "" {
		priceMin, err := strconv.Atoi(priceMinStr)
		if err != nil {
			return nil, err
		}
		params.PriceMin = &priceMin
	}

	if priceMaxStr := query.Get("price_max"); priceMaxStr != "" {
		priceMax, err := strconv.Atoi(priceMaxStr)
		if err != nil {
			return nil, err
		}
		params.PriceMax = &priceMax
	}

	if freePlacesStr := query.Get("free_places"); freePlacesStr != "" {
		freePlaces, err := strconv.Atoi(freePlacesStr)
		if err != nil {
			return nil, err
		}
		params.FreePlaces = &freePlaces
	}

	params.Address = strings.TrimSpace(params.Address)

	if params.OrderBy == "" {
		params.OrderBy = "start_time"
	}

	if params.SortOrder == "" || (params.SortOrder != "asc" && params.SortOrder != "desc") {
		params.SortOrder = "asc"
	}

	return params, nil
}
