package entity

import (
	"fmt"
	"time"
)

type PublicationStatus string

const (
	StatusSent            PublicationStatus = "sent"
	StatusAwaits          PublicationStatus = "awaits"
	StatusErrorOnSending  PublicationStatus = "error_on_sending"
	StatusDeletedByBot    PublicationStatus = "deleted_by_bot"
	StatusErrorOnDeleting PublicationStatus = "error_on_deleting"
)

type Publication struct {
	ID                int               `json:"id"`
	ChannelID         int64             `json:"channel_id"`
	PublicationStatus PublicationStatus `json:"publication_status"`
	Text              string            `json:"text"`
	Image             *string           `json:"image"`
	ButtonUrl         *string           `json:"button_url"`
	ButtonText        *string           `json:"button_text"`
	PublicationDate   *time.Time        `json:"publication_date"`
	DeleteDate        *time.Time        `json:"delete_date"`
	MessageID         int64             `json:"message_id"`

	// channel table - for join
	TelegramChannelID int64  `json:"tg_id"`
	ChannelName       string `json:"channel_name"`
}

func (p Publication) String() string {
	return fmt.Sprintf("(id: %d | channel_id: %d | publication_status: %s | text: %s | image: %v |"+
		" publication_date: %s | delete_date: %v | button_url: %v | button_text: %v)",
		p.ID, p.ChannelID, p.PublicationStatus, p.Text, p.Image, p.PublicationDate, p.DeleteDate, p.ButtonUrl, p.ButtonText)
}
