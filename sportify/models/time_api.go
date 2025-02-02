package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/TheVovchenskiy/sportify-backend/pkg/common"
)

type dateAndTimeAPI struct {
	Date      time.Time `json:"date"`
	StartTime string    `json:"start_time"`
	EndTime   *string   `json:"end_time"`
}

type DateAndTime struct {
	Date      time.Time  `json:"date"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}

func (d *DateAndTime) MarshalJSON() ([]byte, error) {
	startTime := d.StartTime.Format(time.TimeOnly)

	var endTime *string
	if d.EndTime != nil {
		endTime = common.Ref(d.EndTime.Format(time.TimeOnly))
	}

	result := dateAndTimeAPI{
		Date:      d.Date,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return json.Marshal(result)
}

func (d *DateAndTime) UnmarshalJSON(raw []byte) error {
	var dateAndTimeAPI dateAndTimeAPI

	if err := json.Unmarshal(raw, &dateAndTimeAPI); err != nil {
		return err
	}

	year, month, day := dateAndTimeAPI.Date.Date()
	dateAndTimeAPI.Date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	startTime, err := time.Parse(time.TimeOnly, dateAndTimeAPI.StartTime)
	if err != nil {
		return fmt.Errorf("to parse start time: %w", err)
	}

	//startTime = utils.SetTimeZone(dateAndTimeAPI.Date, startTime)

	// TODO later change from UTC to correct
	startTime = time.Date(year, month, day, startTime.Hour(), startTime.Minute(), startTime.Second(), 0, time.UTC)

	var endTime *time.Time
	if dateAndTimeAPI.EndTime != nil {
		endTimeValue, err := time.Parse(time.TimeOnly, *dateAndTimeAPI.EndTime)
		if err != nil {
			return fmt.Errorf("to parse end time: %w", err)
		}

		//endTimeValue = utils.SetTimeZone(dateAndTimeAPI.Date, endTimeValue)
		endTime = &endTimeValue

		// TODO later change from UTC to correct
		endTime = common.Ref(time.Date(year, month, day, endTime.Hour(), endTime.Minute(), endTime.Second(), 0, time.UTC))
	}

	d.Date = dateAndTimeAPI.Date
	d.StartTime = startTime
	d.EndTime = endTime

	return nil
}
