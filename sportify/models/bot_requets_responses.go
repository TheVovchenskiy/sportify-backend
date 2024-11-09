package models

type EventCreatedBotRequest struct {
	TgChatID *int64     `json:"tg_chat_id"`
	TgUserID *int64     `json:"tg_user_id"`
	Event    ShortEvent `json:"event"`
}
