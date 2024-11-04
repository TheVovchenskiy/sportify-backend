package app_test

import (
	"testing"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/app"
	"github.com/TheVovchenskiy/sportify-backend/models"
	"github.com/TheVovchenskiy/sportify-backend/pkg/common"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHashtags(t *testing.T) {
	t.Parallel()

	testSep := " "

	type args struct {
		event *models.ShortEvent
	}

	testCases := map[string]struct {
		args         args
		wantHashtags string
	}{
		"football_free_high_mid_plus": {
			args: args{
				event: &models.ShortEvent{ //nolint:exhaustruct
					SportType:  models.SportTypeFootball,
					Date:       time.Date(2024, 10, 12, 0, 0, 0, 0, time.UTC),
					Price:      nil,
					IsFree:     true,
					GameLevels: []models.GameLevel{models.GameLevelHigh, models.GameLevelMidPlus},
				},
			},
			wantHashtags: "#дата_12_10_2024 #футбол #полу_профи #средний_плюс #бесплатно",
		},
		"basketball_not_free_mid_mid_minus": {
			args: args{
				event: &models.ShortEvent{ //nolint:exhaustruct
					SportType:  models.SportTypeBasketball,
					Date:       time.Date(2024, 10, 24, 0, 0, 0, 0, time.UTC),
					Price:      common.Ref(700),
					IsFree:     false,
					GameLevels: []models.GameLevel{models.GameLevelMid, models.GameLevelMidMinus},
				},
			},
			wantHashtags: "#дата_24_10_2024 #баскетбол #средний #средний_минус #платно",
		},
		"volleyball_not_free_empty_game_levels": {
			args: args{
				event: &models.ShortEvent{ //nolint:exhaustruct
					SportType:  models.SportTypeVolleyball,
					Date:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					Price:      common.Ref(700),
					IsFree:     false,
					GameLevels: []models.GameLevel{},
				},
			},
			wantHashtags: "#дата_01_01_2024 #волейбол #платно",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotHashtags := app.GenerateHashtags(tc.args.event, testSep)

			assert.Equal(t, tc.wantHashtags, gotHashtags)
		})
	}
}
