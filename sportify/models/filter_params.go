package models

import (
	"net/url"
	"strconv"
)

type FilterParams struct {
	SportTypes []string
	GameLevels []string
	DateStarts []string
	PriceMin   *int
	PriceMax   *int
	FreePlaces *int
	OrderBy    string
	SortOrder  string
}

//nolint:cyclop
func ParseFilterParams(query url.Values) (*FilterParams, error) {
	params := &FilterParams{ //nolint:exhaustruct
		SportTypes: query["sport_type"],
		GameLevels: query["game_level"],
		DateStarts: query["date_start"],
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

	if params.OrderBy == "" {
		params.OrderBy = "date_start"
	}

	if params.SortOrder == "" || (params.SortOrder != "asc" && params.SortOrder != "desc") {
		params.SortOrder = "asc"
	}

	return params, nil
}
