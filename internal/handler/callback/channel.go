package callback

import (
	"context"
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/handler"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/tgbot"
	service "github.com/Enthreeka/tg-posting-bot/internal/usecase"
	customErr "github.com/Enthreeka/tg-posting-bot/pkg/bot_error"
	store "github.com/Enthreeka/tg-posting-bot/pkg/local_storage"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
	customMsg "github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api"
	"github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type CallbackChannel interface {
	CallbackShowAllChannels() tgbot.ViewFunc
	CallbackGetChannel() tgbot.ViewFunc
	CallbackCancelCreate() tgbot.ViewFunc
}

type callbackChannel struct {
	channelService     service.ChannelService
	publicationService service.PublicationService
	log                *logger.Logger
	tgMsg              customMsg.Message
	store              store.LocalStorage
}

func NewCallbackChannel(
	channelService service.ChannelService,
	publicationService service.PublicationService,
	log *logger.Logger,
	tgMsg customMsg.Message,
	store store.LocalStorage,
) (CallbackChannel, error) {
	if log == nil {
		return nil, errors.New("logger is nil")
	}
	if channelService == nil {
		return nil, errors.New("userService is nil")
	}
	if tgMsg == nil {
		return nil, errors.New("tgMsg is nil")
	}
	if store == nil {
		return nil, errors.New("storage is nil")
	}
	if publicationService == nil {
		return nil, errors.New("publicationService is nil")
	}

	return &callbackChannel{
		channelService:     channelService,
		publicationService: publicationService,
		log:                log,
		tgMsg:              tgMsg,
		store:              store,
	}, nil
}

// CallbackShowAllChannels - show_channels
func (c *callbackChannel) CallbackShowAllChannels() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelMarkup, err := c.channelService.GetAllAdminChannel(ctx)
		if err != nil {
			c.log.Error("channelService.GetAllAdminChannel: failed to get channel: %v", err)
			handler.HandleError(bot, update, err)
			return nil
		}

		text := `<strong>Ниже представлен список каналов, в которых бот является администратором</strong>`

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			channelMarkup,
			text); err != nil {
			return err
		}

		return nil
	}
}

// CallbackGetChannel - channel_get_{channel_id}/back_setting_{publication_id}
func (c *callbackChannel) CallbackGetChannel() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		var (
			channelIdOrPublicationId = GetID(update.CallbackData())
			channelName              string
		)

		if channelIdOrPublicationId == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}

		if strings.Contains(update.CallbackData(), "channel_get_") {

			channel, err := c.channelService.GetByID(ctx, channelIdOrPublicationId)
			if err != nil {
				c.log.Error("ChannelService.GetByID: failed to get channel: %v", err)
				return err
			}
			channelName = channel.ChannelName
		}
		if strings.Contains(update.CallbackData(), "back_setting_") {

			publication, err := c.publicationService.GetOnePublicationByID(ctx, channelIdOrPublicationId)
			if err != nil {
				c.log.Error("failed to get publication: %v", err)
				return err
			}
			channelIdOrPublicationId = int(publication.ChannelID)
			channelName = publication.ChannelName
		}

		channelSettingMarkup := markup.ChannelSetting(channelIdOrPublicationId)
		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&channelSettingMarkup,
			"Канал: "+channelName); err != nil {
			return err
		}

		return nil
	}
}

// CallbackCancelCreate - cancel-create_{channel_id}
func (c *callbackChannel) CallbackCancelCreate() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		channelID := GetID(update.CallbackData())
		if channelID == 0 {
			c.log.Error("entity.GetID: failed to get id from channel button")
			return customErr.ErrNotFound
		}
		c.store.Delete(update.FromChat().ID)

		channel, err := c.channelService.GetByID(ctx, channelID)
		if err != nil {
			c.log.Error("ChannelService.GetByID: failed to get channel: %v", err)
			return err
		}

		// todo сделать вывод статистики: ожидает отправки

		channelSettingMarkup := markup.ChannelSetting(channelID)
		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&channelSettingMarkup,
			"Канал: "+channel.ChannelName); err != nil {
			return err
		}

		return nil
	}
}
