package models

import "github.com/google/uuid"

type BotUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	TgID     *int64    `json:"tg_id"`
}

type BotEvent struct {
	ID          uuid.UUID   `json:"id"`
	Description *string     `json:"description"`
	Creator     BotUser     `json:"creator"`
	SportType   SportType   `json:"sport_type"`
	Address     string      `json:"address"`
	DateAndTime DateAndTime `json:"date_and_time"`
	Price       *int        `json:"price"`
	IsFree      bool        `json:"is_free"`
	GameLevels  []GameLevel `json:"game_levels"`
	Capacity    *int        `json:"capacity"`
	Busy        int         `json:"busy"`
	Subscribers []BotUser   `json:"subscribers"`
	URLPreview  string      `json:"url_preview"`
	Latitude    *string     `json:"latitude,omitempty"`
	Longitude   *string     `json:"longitude,omitempty"`
	Hashtags    *[]string   `json:"hashtags,omitempty"`
}

type EventCreatedBotRequest struct {
	TgChatID *int64   `json:"tg_chat_id"`
	Event    BotEvent `json:"event"`
	// TgUserID *int64   `json:"tg_user_id"`
}

type EventCreatedBotResponse struct {
	TgChatID    *int64 `json:"tg_chat_id"`
	TgMessageID *int64 `json:"tg_message_id"`
	// TgUserID *int64   `json:"tg_user_id"`
}

type EventUpdatedBotRequest struct {
	TgUserIDsToNotify []int64  `json:"tg_user_ids_to_notify,omitempty"`
	TgChatID          *int64   `json:"tg_chat_id"`
	TgMessageID       *int64   `json:"tg_message_id"`
	Event             BotEvent `json:"event"`
}

type EventDeletedBotRequest struct {
	TgChatID    *int64    `json:"tg_chat_id"`
	TgMessageID *int64    `json:"tg_message_id"`
	EventID     uuid.UUID `json:"event_id"`
}

type SubscribeEventFromTgRequest struct {
	TgChatID    int64 `json:"tg_chat_id"`
	TgMessageID int64 `json:"tg_message_id"`
	TgUserID    int64 `json:"tg_user_id"`
}
