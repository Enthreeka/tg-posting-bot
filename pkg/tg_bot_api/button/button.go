package button

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	MainMenuButton = tgbotapi.NewInlineKeyboardButtonData("Вернуться в главное меню", "main_menu")

	CancelButton = tgbotapi.NewInlineKeyboardButtonData("Отмена выполнения", "cancel")
)
