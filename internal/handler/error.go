package handler

import (
	customErr "github.com/Enthreeka/tg-posting-bot/pkg/bot_error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func HandleError(bot *tgbotapi.BotAPI, update *tgbotapi.Update, err error) {
	msg := tgbotapi.NewMessage(update.FromChat().ID, processError(err))
	if _, err = bot.Send(msg); err != nil {
		log.Printf("failed to send message: %v\n", err)
	}
}

func processError(err error) string {
	if se, ok := err.(*customErr.BotError); ok {
		return se.Msg
	}
	return "Неизвестная ошибка: " + err.Error()
}
