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
	"time"
	"unicode/utf8"
)

type PublicationService interface {
	CreatePublication(ctx context.Context, publication *entity.Publication) (int, error)

	DeletePublication(ctx context.Context, channelId int) error

	GetAllPublicationsByChannelID(ctx context.Context, channelID int, command string) (*tgbotapi.InlineKeyboardMarkup, error)
	GetPublicationByPublicationID(ctx context.Context, publicationID int) (*entity.Publication, error)
	GetAwaitingPublication(ctx context.Context) ([]*entity.Publication, error)
	GetPublicationAndChannel(ctx context.Context, publicationID int) (*entity.Publication, error)
	GetMarkupPublication(publication []entity.Publication, command string) (*tgbotapi.InlineKeyboardMarkup, error)
	GetAllPublicationByID(ctx context.Context, publicationID int) ([]entity.Publication, error)
	GetOnePublicationByID(ctx context.Context, publicationID int) (*entity.Publication, error)

	UpdatePublicationButton(ctx context.Context, publicationID int, buttonUrl, buttonText *string) error
	UpdatePublicationText(ctx context.Context, publicationID int, text string) error
	UpdatePublicationStatus(ctx context.Context, publicationID int, status entity.PublicationStatus) error
	UpdatePublicationImage(ctx context.Context, publicationID int, image *string) error
	UpdatePublicationDate(ctx context.Context, publicationID int, date time.Time) error
	UpdateDeleteDate(ctx context.Context, publicationID int, date time.Time) error
}

type publicationService struct {
	publicationRepo repo.PublicationRepo
	log             *logger.Logger
}

func NewPublicationService(publicationRepo repo.PublicationRepo, log *logger.Logger) (PublicationService, error) {
	if log == nil {
		return nil, errors.New("log is nil")
	}
	if publicationRepo == nil {
		return nil, errors.New("publicationRepo is nil")
	}

	return &publicationService{
		publicationRepo: publicationRepo,
		log:             log,
	}, nil
}

func (p *publicationService) GetOnePublicationByID(ctx context.Context, publicationID int) (*entity.Publication, error) {
	return p.publicationRepo.GetOnePublicationByID(ctx, publicationID)
}

func (p *publicationService) GetAllPublicationByID(ctx context.Context, publicationID int) ([]entity.Publication, error) {
	return p.publicationRepo.GetAllPublicationByID(ctx, publicationID)
}

func (p *publicationService) CreatePublication(ctx context.Context, publication *entity.Publication) (int, error) {
	return p.publicationRepo.CreatePublication(ctx, publication)
}

func (p *publicationService) DeletePublication(ctx context.Context, publicationID int) error {
	return p.publicationRepo.DeletePublication(ctx, publicationID)
}

func (p *publicationService) GetPublicationAndChannel(ctx context.Context, publicationID int) (*entity.Publication, error) {
	publication, err := p.publicationRepo.GetPublicationAndChannel(ctx, publicationID)
	if err != nil {
		return nil, err
	}

	p.log.Info(publication.String())
	return publication, nil
}

func (p *publicationService) UpdatePublicationButton(ctx context.Context, publicationID int, buttonUrl, buttonText *string) error {
	return p.publicationRepo.UpdatePublicationButton(ctx, publicationID, buttonUrl, buttonText)
}

func (p *publicationService) UpdatePublicationText(ctx context.Context, publicationID int, text string) error {
	return p.publicationRepo.UpdatePublicationText(ctx, publicationID, text)
}

func (p *publicationService) UpdatePublicationStatus(ctx context.Context, publicationID int, status entity.PublicationStatus) error {
	return p.publicationRepo.UpdatePublicationStatus(ctx, publicationID, status)
}

func (p *publicationService) UpdatePublicationImage(ctx context.Context, publicationID int, image *string) error {
	return p.publicationRepo.UpdatePublicationImage(ctx, publicationID, image)
}

func (p *publicationService) GetAllPublicationsByChannelID(ctx context.Context, channelID int, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	publication, err := p.publicationRepo.GetAllPublicationByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return p.createPublicationMarkup(publication, command)
}

func (p *publicationService) GetMarkupPublication(publication []entity.Publication, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	return p.createPublicationMarkup(publication, command)
}

func (p *publicationService) UpdatePublicationDate(ctx context.Context, publicationID int, date time.Time) error {
	return p.publicationRepo.UpdatePublicationDate(ctx, publicationID, date)
}

func (p *publicationService) UpdateDeleteDate(ctx context.Context, publicationID int, date time.Time) error {
	return p.publicationRepo.UpdateDeleteDate(ctx, publicationID, date)
}

func (p *publicationService) GetPublicationByPublicationID(ctx context.Context, publicationID int) (*entity.Publication, error) {
	return p.publicationRepo.GetPublicationByPublicationID(ctx, publicationID)
}

func (p *publicationService) GetAwaitingPublication(ctx context.Context) ([]*entity.Publication, error) {
	return p.publicationRepo.GetAwaitingPublication(ctx)
}

func (p *publicationService) createPublicationMarkup(publication []entity.Publication, command string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1
	for i, el := range publication {
		if el.ID != 0 {
			var (
				status string
				text   string
			)
			if el.PublicationStatus == entity.StatusSent {
				status = `✅`
			}
			if el.PublicationStatus == entity.StatusAwaits {
				status = `⏱`
			}
			if el.PublicationStatus == entity.StatusDeletedByBot {
				status = `🗑`
			}

			if el.PublicationStatus == entity.StatusErrorOnSending || el.PublicationStatus == entity.StatusErrorOnDeleting {
				status = `❌`
			}

			if utf8.RuneCountInString(el.Text) < 10 {
				text = el.Text[:len(el.Text)-1]
			} else {
				text = el.Text[:10]
			}

			btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s...%v %s", text, el.PublicationDate.Format(time.DateTime), status),
				fmt.Sprintf("publication_%s_%d", command, el.ID))

			row = append(row, btn)

			if (i+1)%buttonsPerRow == 0 || i == len(publication)-1 {
				rows = append(rows, row)
				row = []tgbotapi.InlineKeyboardButton{}
			}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}
