package models

import (
	"time"

	"github.com/google/uuid"
)

type TgMessage struct {
	RawMessage string `json:"text"`
}

type CreationType string

const (
	CreationTypeTg   CreationType = "tg"
	CreationTypeSite CreationType = "site"
)

const (
	SportTypeVolleyball SportType = "volleyball"
	SportTypeBasketball SportType = "basketball"
	SportTypeFootball   SportType = "football"
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
	URLAuthor    *string
	URLMessage   *string
	CreationType CreationType
	Description  *string `json:"description"`
	RawMessage   *string `json:"raw_message"`
}

type ShortEvent struct {
	ID          uuid.UUID   `json:"id"`
	CreatorID   uuid.UUID   `json:"creator_id"`
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
	URLPreview  string      `json:"preview"`
	URLPhotos   []string    `json:"photos"`
}
