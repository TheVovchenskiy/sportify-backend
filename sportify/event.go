package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//“id”: “uuid (бэк генерит)”,
//"sport_type": "enum(volleball, basketball, football)",
//"address": "text",
//"date" : "timestamp",
//"start_time": "timestamp",
//"end_time": "?timestamp",
//"price": "? number rubles",
//"is_free": "bool",
//"game_level": "? [enums] (low, mid_minus, mid, mid_plus, high",
//"description": "? text",
//"raw_message": "? text",
//"capacity": "? number",
//"busy": "? number",
//“subscribers_id”: “[uuid]”
//“preview”: “string url”,
//“photos”: “[string url]”

type TgMessage struct {
	RawMessage string `json:"text"`
}

const (
	TypeVolleyball SportType = "volleyball"
	TypeBasketball SportType = "basketball"
	TypeFootball   SportType = "football"
)

type SportType string

const (
	GameLevelLow      GameLevel = "low"
	GameLevelMidMinus GameLevel = "mid_minus"
	GameLevelMid      GameLevel = "mid"
	GameLevelMidPlus  GameLevel = "mid_plus"
	GameLevelHigh     GameLevel = "high"
)

type GameLevel string

type FullEvent struct {
	ShortEvent
	Description *string `json:"description"`
	RawMessage  *string `json:"raw_message"`
}

type ShortEvent struct {
	ID          uuid.UUID   `json:"id"`
	SportType   SportType   `json:"sport_type"`
	Address     string      `json:"address"`
	Date        time.Time   `json:"date"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     *time.Time  `json:"end_time"`
	Price       *int        `json:"price"`
	IsFree      bool        `json:"is_free"`
	GameLevel   *GameLevel  `json:"game_level"`
	Capacity    *int        `json:"capacity"`
	Busy        int         `json:"busy"`
	Subscribers []uuid.UUID `json:"subscribers_id"`
	PreviewURL  string      `json:"preview"`
	PhotoURLs   []string    `json:"photos"`
}

var (
	ErrAllBusy            = errors.New("all places are busy")
	ErrNotFoundSubscriber = errors.New("not found subscriber in event")
)

func (s *ShortEvent) AddSubscriber(id uuid.UUID) ([]uuid.UUID, error) {
	if s.Capacity != nil && *s.Capacity <= s.Busy {
		return s.Subscribers, ErrAllBusy
	}

	s.Subscribers = append(s.Subscribers, id)
	s.Busy = len(s.Subscribers)

	return s.Subscribers, nil
}

func (s *ShortEvent) RemoveSubscriber(id uuid.UUID) ([]uuid.UUID, error) {
	for i, v := range s.Subscribers {
		if v == id {
			s.Subscribers = append(s.Subscribers[:i], s.Subscribers[i+1:]...)

			s.Busy = len(s.Subscribers)

			return s.Subscribers, nil
		}
	}

	return s.Subscribers, ErrNotFoundSubscriber
}
