package tgbot

import (
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func userUpdateToModel(update *tgbotapi.Update) *entity.User {
	user := new(entity.User)

	if update != nil {
		user.ID = update.Message.From.ID
		user.TGUsername = update.Message.From.UserName
		user.CreatedAt = time.Now().Local()
		user.UserRole = entity.UserType
	}

	return user
}

func channelUpdateToModel(update *tgbotapi.Update) *entity.Channel {
	channel := &entity.Channel{
		TgID:          update.MyChatMember.Chat.ID,
		ChannelName:   update.MyChatMember.Chat.Title,
		ChannelStatus: entity.GetChannelStatus(update.MyChatMember.NewChatMember.Status),
	}

	if update.MyChatMember.Chat.UserName != "" {
		url := "t.me/" + update.MyChatMember.Chat.UserName
		channel.ChannelUrl = &url
	}

	return channel
}
