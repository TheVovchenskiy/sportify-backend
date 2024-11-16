package utils_test

import (
	"testing"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestSetTimeZone(t *testing.T) {
	t.Parallel()

	type args struct {
		timeZoneFrom time.Time
		timeZoneTo   time.Time
	}

	zoneMSK := time.FixedZone("UTC-8", 3*3600)
	zoneVladivostok := time.FixedZone("UTC-8", 10*3600)
	zoneUTCMinusOne := time.FixedZone("UTC-8", -1*3600)

	testCases := map[string]struct {
		args           args
		wantTimeZoneTo time.Time
	}{
		"UTC_to_Moscow": {args: args{
			timeZoneFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, zoneMSK),
			timeZoneTo:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
			wantTimeZoneTo: time.Date(2000, 1, 1, 0, 0, 0, 0, zoneMSK),
		},
		"Moscow_to_Moscow": {args: args{
			timeZoneFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, zoneMSK),
			timeZoneTo:   time.Date(2000, 1, 1, 0, 0, 0, 0, zoneMSK),
		},
			wantTimeZoneTo: time.Date(2000, 1, 1, 0, 0, 0, 0, zoneMSK),
		},
		"Moscow_to_zoneVladivostok": {args: args{
			timeZoneFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, zoneVladivostok),
			timeZoneTo:   time.Date(2012, 1, 1, 7, 0, 0, 0, zoneMSK),
		},
			wantTimeZoneTo: time.Date(2012, 1, 1, 7, 0, 0, 0, zoneVladivostok),
		},
		"UTC-1_to_zoneVladivostok": {args: args{
			timeZoneFrom: time.Date(2020, 1, 1, 0, 0, 0, 0, zoneVladivostok),
			timeZoneTo:   time.Date(2012, 1, 1, 7, 0, 0, 0, zoneUTCMinusOne),
		},
			wantTimeZoneTo: time.Date(2012, 1, 1, 7, 0, 0, 0, zoneVladivostok),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			gotTime := utils.SetTimeZone(tc.args.timeZoneFrom, tc.args.timeZoneTo)

			assert.Equal(t, tc.wantTimeZoneTo, gotTime)
		})
	}
}
