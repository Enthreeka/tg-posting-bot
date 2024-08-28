package entity

import (
	"fmt"
)

type ChannelStatus string

const (
	StatusUnknown       ChannelStatus = "unknown"
	StatusKicked        ChannelStatus = "kicked"
	StatusAdministrator ChannelStatus = "administrator"
	StatusLeft          ChannelStatus = "left"
	StatusMember        ChannelStatus = "member"
)

func GetChannelStatus(gotStatus string) ChannelStatus {
	switch gotStatus {
	case "kicked":
		return StatusKicked
	case "administrator":
		return StatusAdministrator
	case "left":
		return StatusLeft
	case "member":
		return StatusMember
	default:
		return StatusUnknown
	}
}

type Channel struct {
	ID            int           `json:"id"`
	TgID          int64         `json:"tg_id"`
	ChannelName   string        `json:"channel_name"`
	ChannelUrl    *string       `json:"channel_url"`
	ChannelStatus ChannelStatus `json:"channel_status"`
}

func (c Channel) String() string {
	var url string
	if c.ChannelUrl == nil {
		url = "none"
	} else {
		url = *c.ChannelUrl
	}

	return fmt.Sprintf("(tg_id: %d | channel_name: %s | ChannelURL: %s | Status: %s)",
		c.TgID, c.ChannelName, url, c.ChannelStatus)
}
