package models

import "fmt"

type TgMessage struct {
	MessageID int `json:"message_id"`
	Chat      struct {
		Username string `json:"username"`
	} `json:"chat"`
	RawMessage string `json:"text"`
	From       struct {
		Username string `json:"username"`
	} `json:"from"`
}

const (
	formatURLAuthor  = "https://t.me/%s"
	formatURLMessage = "https://t.me/%s/%d"
)

func (t *TgMessage) GetURLAuthor() string {
	return fmt.Sprintf(formatURLAuthor, t.From.Username)
}

func (t *TgMessage) GetURLMessage() string {
	return fmt.Sprintf(formatURLMessage, t.Chat.Username, t.MessageID)
}
