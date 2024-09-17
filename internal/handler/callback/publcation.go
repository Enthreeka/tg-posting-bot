package callback

import (
	"context"
	"errors"
	"fmt"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/tgbot"
	service "github.com/Enthreeka/tg-posting-bot/internal/usecase"
	customErr "github.com/Enthreeka/tg-posting-bot/pkg/bot_error"
	store "github.com/Enthreeka/tg-posting-bot/pkg/local_storage"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
	customMsg "github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api"
	"github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PublicationChannel interface {
	CallbackCreatePublication() tgbot.ViewFunc
	CallbackUpdatePublicationSettings() tgbot.ViewFunc
	CallbackUpdatePublicationText() tgbot.ViewFunc
	CallbackUpdatePublicationImage() tgbot.ViewFunc
	CallbackUpdatePublicationButton() tgbot.ViewFunc
	CallbackUpdatePublicationSentDate() tgbot.ViewFunc
	CallbackUpdatePublicationDeleteDate() tgbot.ViewFunc
	CallbackCheckPublication() tgbot.ViewFunc
	CallbackGetPublicationGet() tgbot.ViewFunc
	CallbackGetListForCancelPublication() tgbot.ViewFunc
	CallbackDeletePublication() tgbot.ViewFunc
	CallbackCancelUpdate() tgbot.ViewFunc
}

type callbackPublication struct {
	publicationService service.PublicationService
	channelService     service.ChannelService
	log                *logger.Logger
	tgMsg              customMsg.Message
	store              store.LocalStorage
	publicationArray   *store.PublicationArray
}

func NewCallbackPublication(
	publicationService service.PublicationService,
	log *logger.Logger,
	tgMsg customMsg.Message,
	store store.LocalStorage,
	channelService service.ChannelService,
	publicationArray *store.PublicationArray,
) (PublicationChannel, error) {
	if log == nil {
		return nil, errors.New("logger is nil")
	}
	if publicationService == nil {
		return nil, errors.New("publicationService is nil")
	}
	if tgMsg == nil {
		return nil, errors.New("tgMsg is nil")
	}
	if store == nil {
		return nil, errors.New("store is nil")
	}
	if channelService == nil {
		return nil, errors.New("channelService is nil")
	}
	if publicationArray == nil {
		return nil, errors.New("publicationArray is nil")
	}

	return &callbackPublication{
		publicationService: publicationService,
		channelService:     channelService,
		log:                log,
		tgMsg:              tgMsg,
		store:              store,
		publicationArray:   publicationArray,
	}, nil
}

// CallbackCreatePublication - create_publication
func (c *callbackPublication) CallbackCreatePublication() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelID := GetID(update.CallbackData())
		if channelID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		text := "Для создания публикации отправьте сообщение следующего вида:\n" +
			`	{
			  "дата_публикации": "2024-08-27 15:48",
			  "дата_удаления": "2024-08-27 15:49",
			  "кнопка": {
				"текст_кнопки": "текст кнопки",
				"ссылка_кнопки": "https://yandex.ru"
			  }
			}` +
			"\n\n Обязательными полями являются: дата_публикации\n\n" +
			`{
				 "дата_публикации": "2024-08-27 15:48"
				}`

		cancelCommandMarkup := markup.CancelCommandCreate(channelID)
		sentMsg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&cancelCommandMarkup,
			text)
		if err != nil {
			return err
		}

		c.store.Set(&store.Data{
			CurrentMsgID:  sentMsg,
			PreferMsgID:   update.CallbackQuery.Message.MessageID,
			OperationType: store.PublicationCreate,
			ChannelID:     channelID,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackGetPublicationList - publication_get_{publication_id}
func (c *callbackPublication) CallbackGetPublicationGet() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		publication, err := c.publicationService.GetPublicationAndChannel(ctx, publicationID)
		if err != nil {
			c.log.Error("failed to GetPublicationAndChannel: %v", err)
			return err
		}

		text := fmt.Sprintf("*Изменение публикации*\n\n"+
			"Канал: %s\n"+
			"Время удаления: %v\n"+
			"Время отправления: %v", publication.ChannelName, publication.DeleteDate, publication.PublicationDate)
		updatePublicationSettingsMarkup := markup.UpdatePublicationSettings(publicationID)
		_, err = c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&updatePublicationSettingsMarkup,
			text)
		if err != nil {
			return err
		}

		return nil
	}
}

// CallbackUpdatePublicationSettings - publication_update_{channel_id}
func (c *callbackPublication) CallbackUpdatePublicationSettings() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelID := GetID(update.CallbackData())
		if channelID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		channel, err := c.channelService.GetByID(ctx, channelID)
		if err != nil {
			return err
		}

		publicationMarkup, err := c.publicationService.GetAllPublicationsByChannelID(ctx, channelID, "get")
		if err != nil {
			return err
		}

		text := fmt.Sprintf("Публикации для канала: **%s**\n\n❌ - ошибка при удалении/ошибка при отправке\n"+
			"✅ - отправлено\n⏱ - ожидает отправки\n🗑 - удалено из канала", channel.ChannelName)
		_, err = c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			publicationMarkup,
			text)
		if err != nil {
			return err
		}

		return nil
	}
}

// CallbackUpdatePublicationText - text_update_{publication_id}
func (c *callbackPublication) CallbackUpdatePublicationText() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		text := "Отправьте новый текст публикации"
		cancelCommandMarkup := markup.CancelCommandPublication(publicationID)
		sentMsg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&cancelCommandMarkup,
			text)
		if err != nil {
			return err
		}

		c.store.Set(&store.Data{
			CurrentMsgID:  sentMsg,
			PreferMsgID:   update.CallbackQuery.Message.MessageID,
			OperationType: store.PublicationTextUpdate,
			ChannelID:     publicationID,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackUpdatePublicationImage - image_update_{publication_id}
func (c *callbackPublication) CallbackUpdatePublicationImage() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		text := "Отправьте изображение для публикации"
		cancelCommandMarkup := markup.CancelCommandPublication(publicationID)
		sentMsg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&cancelCommandMarkup,
			text)
		if err != nil {
			return err
		}

		c.store.Set(&store.Data{
			CurrentMsgID:  sentMsg,
			PreferMsgID:   update.CallbackQuery.Message.MessageID,
			OperationType: store.PublicationImageUpdate,
			ChannelID:     publicationID,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackUpdatePublicationButton - button_update_{publication_id}
func (c *callbackPublication) CallbackUpdatePublicationButton() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		text := "Отправьте подпись для кнопки и ссылку (ссылка не является обязательной).\n" +
			"Пример с ссылкой:\n" +
			`{
				"текст_кнопки": "текст кнопки",
				"ссылка_кнопки": "https://yandex.ru"
			}` +
			"Пример без ссылки:\n" +
			`{
				"текст_кнопки": "текст кнопки"
			}`
		cancelCommandMarkup := markup.CancelCommandPublication(publicationID)
		sentMsg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&cancelCommandMarkup,
			text)
		if err != nil {
			return err
		}

		c.store.Set(&store.Data{
			CurrentMsgID:  sentMsg,
			PreferMsgID:   update.CallbackQuery.Message.MessageID,
			OperationType: store.PublicationButtonUpdate,
			ChannelID:     publicationID,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackUpdatePublicationSentDate - sent-date_update_{publication_id}
func (c *callbackPublication) CallbackUpdatePublicationSentDate() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		text := "Отправьте время и дату в формате: 2024-08-27 15:48"
		cancelCommandMarkup := markup.CancelCommandPublication(publicationID)
		sentMsg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&cancelCommandMarkup,
			text)
		if err != nil {
			return err
		}

		c.store.Set(&store.Data{
			CurrentMsgID:  sentMsg,
			PreferMsgID:   update.CallbackQuery.Message.MessageID,
			OperationType: store.PublicationSentDateUpdate,
			ChannelID:     publicationID,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackUpdatePublicationDeleteDate - delete-date_update_{publication_id}
func (c *callbackPublication) CallbackUpdatePublicationDeleteDate() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		text := "Отправьте время и дату в формате: 2024-08-27 15:48"
		cancelCommandMarkup := markup.CancelCommandPublication(publicationID)
		sentMsg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&cancelCommandMarkup,
			text)
		if err != nil {
			return err
		}

		c.store.Set(&store.Data{
			CurrentMsgID:  sentMsg,
			PreferMsgID:   update.CallbackQuery.Message.MessageID,
			OperationType: store.PublicationDeleteDateUpdate,
			ChannelID:     publicationID,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackCheckPublication - check_publication_{publication_id}
func (c *callbackPublication) CallbackCheckPublication() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		publication, err := c.publicationService.GetPublicationByPublicationID(ctx, publicationID)
		if err != nil {
			return err
		}

		if _, err := c.tgMsg.SendMessageToUser(update.FromChat().ID, publication); err != nil {
			c.log.Error("failed to send message to user: %v", err)
			return err
		}

		return nil
	}
}

// CallbackGetListForCancelPublication - publication_cancel
func (c *callbackPublication) CallbackGetListForCancelPublication() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelID := GetID(update.CallbackData())
		if channelID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		channel, err := c.channelService.GetByID(ctx, channelID)
		if err != nil {
			return err
		}

		publicationMarkup, err := c.publicationService.GetAllPublicationsByChannelID(ctx, channelID, "delete")
		if err != nil {
			return err
		}

		text := fmt.Sprintf("Публикации для канала: **%s**\\n\\n❌ - ошибка при удалении|ошибка при отправке\\n\"+\n\t\t\t\"✅ - отправлено\\n⏱ - ожидает отправки\\n🗑 - удалено из канала"+
			"\n\n Нажмите на публикацию чтобы ее *удалить*", channel.ChannelName)
		_, err = c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			publicationMarkup,
			text)
		if err != nil {
			return err
		}

		return nil
	}
}

// CallbackDeletePublication - publication_delete_{publication_id}
func (c *callbackPublication) CallbackDeletePublication() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		publication, err := c.publicationService.GetPublicationByPublicationID(ctx, publicationID)
		if err != nil {
			return err
		}

		if err := c.publicationService.DeletePublication(ctx, publicationID); err != nil {
			c.log.Error("failed to delete publication: %v", err)
		}

		c.publicationArray.RemovePub(&store.PubData{
			PubDate:       publication.PublicationDate,
			PublicationID: publicationID,
		})

		text := "Публикация удалена"
		if _, err := c.tgMsg.SendNewMessage(update.FromChat().ID, nil, text); err != nil {
			return err
		}

		return nil
	}
}

func (c *callbackPublication) CallbackCancelUpdate() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		publicationID := GetID(update.CallbackData())
		if publicationID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		c.store.Delete(update.FromChat().ID)

		publication, err := c.publicationService.GetAllPublicationByID(ctx, publicationID)
		if err != nil {
			c.log.Error("failed to get publication: %v", err)
			return err
		}

		publicationMarkup, err := c.publicationService.GetMarkupPublication(publication, "get")
		if err != nil {
			return err
		}

		text := fmt.Sprintf("Публикации для канала: **%s**\n\n❌ - ошибка при удалении|ошибка при отправке\n"+
			"✅ - отправлено\n⏱ - ожидает отправки\n🗑 - удалено из канала", publication[0].ChannelName)
		_, err = c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			publicationMarkup,
			text)
		if err != nil {
			return err
		}

		return nil
	}
}
