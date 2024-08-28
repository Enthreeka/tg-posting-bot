package markup

import (
	"fmt"
	"github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	StartMenu = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление ботом", "show_channels")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление пользователями", "user_setting")),
	)

	UserSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить роль администратора", "admin_set_role"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отозвать роль администратора", "admin_delete_role"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Посмотреть список администраторов", "admin_look_up"),
		),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	SuperAdminSetting = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить администратором", "create_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назначить супер администратором", "create_super_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Забрать права администратора", "delete_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Список администраторов", "all_admin")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "user_setting")),
		tgbotapi.NewInlineKeyboardRow(button.MainMenuButton),
	)

	MainMenu = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button.MainMenuButton))

	//CancelCommand = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button.CancelButton))
)

func CancelCommandCreate(channelID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена выполнения", fmt.Sprintf("cancel_create_%d", channelID))))
}

func CancelCommandPublication(publicationId int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отмена выполнения", fmt.Sprintf("cancel_update_%d", publicationId))))
}

func ChannelSetting(channelID int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Создать публикацию", fmt.Sprintf("publication_create_%d", channelID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Управление публикациями", fmt.Sprintf("publication_update_%d", channelID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отменить публикацию", fmt.Sprintf("publication_cancel_%d", channelID))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", "show_channels")),
	)
}

func UpdatePublicationSettings(publicationId int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить описание", fmt.Sprintf("text_update_%d", publicationId))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить фотографию", fmt.Sprintf("image_update_%d", publicationId))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить кнопку", fmt.Sprintf("button_update_%d", publicationId))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить дату отправки", fmt.Sprintf("sent-date_update_%d", publicationId))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить дату удаления", fmt.Sprintf("delete-date_update_%d", publicationId))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предварительный просмотр", fmt.Sprintf("check_publication_%d", publicationId))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вернуться назад", fmt.Sprintf("back_setting_%d", publicationId))),
	)
}
