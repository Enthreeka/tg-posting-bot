package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	"github.com/Enthreeka/tg-posting-bot/internal/repo"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
	"github.com/Enthreeka/tg-posting-bot/pkg/tg_bot_api/button"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChannelService interface {
	Create(ctx context.Context, channel *entity.Channel) error

	GetByID(ctx context.Context, id int) (*entity.Channel, error)
	GetAll(ctx context.Context) ([]entity.Channel, error)
	GetAllAdminChannel(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error)
	GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error)

	DeleteByID(ctx context.Context, id int) error
	ChatMember(ctx context.Context, channel *entity.Channel) error
}

type channelService struct {
	channelRepo repo.ChannelRepo
	log         *logger.Logger
}

func NewChannelService(channelRepo repo.ChannelRepo, log *logger.Logger) (ChannelService, error) {
	if log == nil {
		return nil, errors.New("log is nil")
	}
	if channelRepo == nil {
		return nil, errors.New("channelRepo is nil")
	}

	return &channelService{
		channelRepo: channelRepo,
		log:         log,
	}, nil
}

func (c *channelService) Create(ctx context.Context, channel *entity.Channel) error {
	return c.channelRepo.Create(ctx, channel)
}

func (c *channelService) GetByID(ctx context.Context, id int) (*entity.Channel, error) {
	return c.channelRepo.GetByID(ctx, id)
}

func (c *channelService) DeleteByID(ctx context.Context, id int) error {
	return c.channelRepo.DeleteByID(ctx, id)
}

func (c *channelService) GetAll(ctx context.Context) ([]entity.Channel, error) {
	return c.channelRepo.GetAll(ctx)
}

func (c *channelService) ChatMember(ctx context.Context, channel *entity.Channel) error {
	c.log.Info("GetPub channel: %s", channel.String())

	isExist, err := c.channelRepo.IsChannelExistByTgID(ctx, channel.TgID)
	if err != nil {
		c.log.Error("channelRepo.IsChannelExistByTgID: failed to check channel: %v", err)
		return err
	}

	if !isExist {
		err := c.channelRepo.Create(ctx, channel)
		if err != nil {
			c.log.Error("channelRepo.Create: failed to create channel: %v", err)
			return err
		}
		return nil
	}

	err = c.channelRepo.UpdateStatusByTgID(ctx, channel.ChannelStatus, channel.TgID)
	if err != nil {
		c.log.Error("channelRepo.UpdateStatusByTgID: failed to update channel status: %v", err)
		return err
	}
	return nil
}

func (c *channelService) GetAllAdminChannel(ctx context.Context) (*tgbotapi.InlineKeyboardMarkup, error) {
	channel, err := c.channelRepo.GetAllAdminChannel(ctx)
	if err != nil {
		return nil, err
	}

	return c.createChannelMarkup(channel, "get")
}

func (c *channelService) createChannelMarkup(channel []entity.Channel, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1
	for i, el := range channel {
		if el.TgID != 0 { // check for channel for global notification
			btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.ChannelName),
				fmt.Sprintf("channel_%s_%d", command, el.ID))

			row = append(row, btn)

			if (i+1)%buttonsPerRow == 0 || i == len(channel)-1 {
				rows = append(rows, row)
				row = []tgbotapi.InlineKeyboardButton{}
			}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (c *channelService) GetByChannelName(ctx context.Context, channelName string) (*entity.Channel, error) {
	return c.channelRepo.GetByChannelName(ctx, channelName)
}
