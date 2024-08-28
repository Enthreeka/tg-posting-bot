package tgbot

import (
	store "github.com/Enthreeka/tg-posting-bot/pkg/local_storage"
	"github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	success = "Операция выполнена успешно. "
)

// response - возвращает ответ администратору
func (b *Bot) response(operationType store.TypeCommand, currentMessageId int, preferMessageId int, channelID int, update *tgbotapi.Update) {
	var (
		messageId int
		userID    = update.FromChat().ID
	)

	if update.Message != nil {
		messageId = update.Message.MessageID
	} else if update.CallbackQuery != nil {
		messageId = update.CallbackQuery.Message.MessageID
	}

	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID, messageId)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", messageId, string(resp.Result), err)
	}

	// Выполнять удаление сообщения только для определенных операций
	if value, _ := store.MapTypes[operationType]; value == store.Admin {
		if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID, currentMessageId)); nil != err || !resp.Ok {
			b.log.Error("failed to delete message id %d (%s): %v", currentMessageId, string(resp.Result), err)
		}
	}

	text, markup := responseText(operationType, channelID)
	if _, err := b.tgMsg.SendEditMessage(userID, preferMessageId, markup, text); err != nil {
		b.log.Error("failed to send telegram message: ", err)
	}
}

func responseText(operationType store.TypeCommand, channelID int) (string, *tgbotapi.InlineKeyboardMarkup) {
	switch operationType {
	case store.AdminCreate:
		return success + "Пользователь получил администраторские права.", &markup.UserSetting
	case store.AdminDelete:
		return success + "Пользователь лишился администраторских прав.", &markup.UserSetting
	case store.PublicationCreate:
		keyMarkup := markup.ChannelSetting(channelID)
		return success + "Публикация добавлена.", &keyMarkup
	}
	return success, nil
}
