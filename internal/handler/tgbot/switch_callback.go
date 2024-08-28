package tgbot

import (
	customErr "github.com/Enthreeka/tg-posting-bot/pkg/bot_error"
	"strings"
)

func (b *Bot) CallbackStrings(callbackData string) (error, ViewFunc) {
	for key, _ := range b.callbackStore.GetStorage() {
		if strings.HasPrefix(callbackData, key) || strings.HasPrefix(callbackData, key+"_") {
			callbackView, ok := b.callbackView[key]
			if !ok {
				return customErr.ErrNotFound, nil
			}
			return nil, callbackView
		}
	}
	return nil, nil
}
