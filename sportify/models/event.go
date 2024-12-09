package models

import (
	"github.com/google/uuid"
)

type CreationType string

const (
	CreationTypeTg   CreationType = "tg"
	CreationTypeSite CreationType = "site"
)

type CreatorAPI struct {
	CreatorID uuid.UUID `json:"creator_id"`
	Username  string    `json:"username"`
	TgURL     *string   `json:"tg_url"`
}

type FullEventAPI struct {
	FullEvent
	CreatorAPI CreatorAPI `json:"creator"`
}

func MapFullEventToAPI(fullEvent *FullEvent, username string, tgUserID *int64) *FullEventAPI {
	return &FullEventAPI{
		FullEvent: *fullEvent,
		CreatorAPI: CreatorAPI{
			CreatorID: fullEvent.CreatorID,
			Username:  username,
			TgURL:     mapTgURL(tgUserID, username),
		},
	}
}

type FullEvent struct {
	ShortEvent
	URLAuthor    *string      `json:"url_author"`
	URLMessage   *string      `json:"url_message"`
	CreationType CreationType `json:"creation_type"`
	Description  *string      `json:"description"`
	RawMessage   *string      `json:"raw_message"`
}

func NewFullEventSite(eventID uuid.UUID, userID uuid.UUID, eventCreteSite *EventCreateSite) *FullEvent {
	return &FullEvent{ //nolint:exhaustruct
		ShortEvent: ShortEvent{
			ID:        eventID,
			CreatorID: userID,
			SportType: eventCreteSite.SportType,
			Address:   eventCreteSite.Address,
			DateAndTime: DateAndTime{
				Date:      eventCreteSite.DateAndTime.Date,
				StartTime: eventCreteSite.DateAndTime.StartTime,
				EndTime:   eventCreteSite.DateAndTime.EndTime,
			},
			Price:       eventCreteSite.Price,
			IsFree:      IsFreePrice(eventCreteSite.Price),
			GameLevels:  eventCreteSite.GameLevels,
			Capacity:    eventCreteSite.Capacity,
			Busy:        0,
			Subscribers: make([]uuid.UUID, 0),
			URLPreview:  eventCreteSite.URLPreview,
			URLPhotos:   eventCreteSite.URLPhotos,
		},
		CreationType: CreationTypeSite,
		Description:  eventCreteSite.Description,
	}
}

func (e *FullEvent) ToBotEvent(creator *BotUser, subscribers []*BotUser, hashtags *[]string) *BotEvent {
	subs := make([]BotUser, 0)
	for _, sub := range subscribers {
		subs = append(subs, *sub)
	}

	return &BotEvent{ //nolint:exhaustruct
		ID:          e.ID,
		Creator:     *creator,
		SportType:   e.SportType,
		Address:     e.Address,
		DateAndTime: e.DateAndTime,
		Price:       e.Price,
		IsFree:      e.IsFree,
		GameLevels:  e.GameLevels,
		Capacity:    e.Capacity,
		Busy:        e.Busy,
		Subscribers: subs,
		URLPreview:  e.URLPreview,
		Latitude:    e.Latitude,
		Longitude:   e.Longitude,
		Hashtags:    hashtags,
	}
}

type ShortEvent struct {
	ID          uuid.UUID   `json:"id"`
	CreatorID   uuid.UUID   `json:"creator_id"`
	SportType   SportType   `json:"sport_type"`
	Address     string      `json:"address"`
	DateAndTime DateAndTime `json:"date_time"`
	Price       *int        `json:"price"`
	IsFree      bool        `json:"is_free"`
	GameLevels  []GameLevel `json:"game_level"`
	Capacity    *int        `json:"capacity"`
	Busy        int         `json:"busy"`
	Subscribers []uuid.UUID `json:"subscribers_id"`
	URLPreview  string      `json:"preview"`
	URLPhotos   []string    `json:"photos"`
	Latitude    *string     `json:"latitude"`
	Longitude   *string     `json:"longitude"`
}

func IsFreePrice(price *int) bool {
	return price == nil || *price == 0
}

func RawGameLevel(gameLevels []GameLevel) []string {
	result := make([]string, len(gameLevels))

	for i, v := range gameLevels {
		result[i] = string(v)
	}

	return result
}

func GameLevelFromRaw(gameLevels []string) []GameLevel {
	result := make([]GameLevel, len(gameLevels))

	for i, v := range gameLevels {
		result[i] = GameLevel(v)
	}

	return result
}

func GameLevelFromRawNullable(gameLevels []*string) []GameLevel {
	result := []GameLevel{} // for contract with front

	for _, v := range gameLevels {
		if v == nil {
			continue
		}
		result = append(result, GameLevel(*v))
	}

	return result
}
