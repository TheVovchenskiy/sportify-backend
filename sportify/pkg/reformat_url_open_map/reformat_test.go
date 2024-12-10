package reformat_url_open_map_test

import (
	"github.com/TheVovchenskiy/sportify-backend/pkg/reformat_url_open_map"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReformatURLOpenMap(t *testing.T) {
	t.Parallel()

	testCases := map[string]string{
		"Госпитальный переулок 4-6":               "Госпитальный переулок 4-6",
		"Алтайский край":                          "Алтайский край",
		"г Москва, Госпитальная наб, д 4 стр 1":   "Москва, Госпитальная набережная, 4",
		"г Москва, ул Воротынская, д 9 к 1":       "Москва, улица Воротынская, 9",
		"г Москва, Госпитальный пер, д 4-6 стр 3": "Москва, Госпитальный переулок, 4-6",
		"Московская обл, г Клин, деревня Кононово, тер. СНТ Аллея Перова МГТУ им Н.Э.Баумана": "Клин, деревня Кононово, территория СНТ Аллея Перова МГТУ им Н.Э.Баумана",
	}

	for input, want := range testCases {
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			got := reformat_url_open_map.ReformatURLOpenMap(input)
			assert.Equal(t, want, got)
		})
	}
}
