package tgbot

import (
	"context"
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/tgbot/dto"
	store "github.com/Enthreeka/tg-posting-bot/pkg/local_storage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/url"
	"time"
	"unicode/utf16"
)

var needEscape = make(map[rune]struct{})

func init() {
	for _, r := range []rune{'_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'} {
		needEscape[r] = struct{}{}
	}
}

func (b *Bot) isStateExist(userID int64) (*store.Data, bool) {
	data, exist := b.store.Read(userID)
	return data, exist
}

func (b *Bot) isStoreProcessing(ctx context.Context, update *tgbotapi.Update) (bool, error) {
	userID := update.Message.From.ID
	storeData, isExist := b.isStateExist(userID)
	if !isExist || storeData == nil {
		return false, nil
	}
	defer b.store.Delete(userID)

	return b.switchStoreData(ctx, update, storeData)
}

// todo refactor channelID where real it is publicationID
func (b *Bot) switchStoreData(ctx context.Context, update *tgbotapi.Update, storeData *store.Data) (bool, error) {
	var (
		err error
	)

	switch storeData.OperationType {
	case store.AdminCreate:
		if err = b.userService.UpdateRoleByUsername(ctx, entity.AdminType, update.Message.Text); err != nil {
			b.log.Error("isStoreExist::store.AdminCreate:UpdateRoleByUsername: %v", err)
		}
	case store.AdminDelete:
		if err = b.userService.UpdateRoleByUsername(ctx, entity.UserType, update.Message.Text); err != nil {
			b.log.Error("isStoreExist::store.AdminDelete:userRepo.UpdateRoleByUsername: %v", err)
		}
	case store.PublicationCreate:

		_, err = b.publicationService.CreatePublicationOnlyWithText(ctx, ConvertToMarkdownV2(update.Message.Text, update.Message.Entities), storeData.ChannelID)
		if err != nil {
			b.log.Error("isStoreExist::store.PublicationCreate: %v", err)
		}
	case store.PublicationTextUpdate:
		// todo переделать с storeData.ChannelID на storeData.PublicationID
		if err = b.publicationService.UpdatePublicationText(ctx, storeData.ChannelID, ConvertToMarkdownV2(update.Message.Text, update.Message.Entities)); err != nil {
			b.log.Error("isStoreExist::store.PublicationTextUpdate: %v", err)
		}
	case store.PublicationImageUpdate:
		largestPhoto := update.Message.Photo[len(update.Message.Photo)-1]
		// todo переделать с storeData.ChannelID на storeData.PublicationID
		if err = b.publicationService.UpdatePublicationImage(ctx, storeData.ChannelID, &largestPhoto.FileID); err != nil {
			b.log.Error("isStoreExist::store.PublicationTextImage: %v", err)
		}
	case store.PublicationButtonTextUpdate:
		// todo переделать с storeData.ChannelID на storeData.PublicationID
		if err = b.publicationService.UpdatePublicationButtonText(ctx, storeData.ChannelID, update.Message.Text); err != nil {
			b.log.Error("isStoreExist::store.PublicationButtonTextUpdate: %v", err)
		}
	case store.PublicationButtonLinkUpdate:
		// todo переделать с storeData.ChannelID на storeData.PublicationID
		if _, err = url.ParseRequestURI(update.Message.Text); err != nil {
			b.log.Error("isStoreExist::store.PublicationButtonLinkUpdate: %v", err)
			return true, errors.New("ошибка: невалидная ссылка")
		}

		if err = b.publicationService.UpdatePublicationButtonLink(ctx, storeData.ChannelID, update.Message.Text); err != nil {
			b.log.Error("isStoreExist::store.PublicationButtonLinkUpdate: %v", err)
		}
	case store.PublicationDeleteDateUpdate:
		var (
			date time.Time
		)
		date, err = time.Parse(dto.Layout, update.Message.Text+":00 +0300")
		if err != nil {
			b.log.Error("isStoreExist::store.PublicationDeleteDateUpdate: %v", err)
			return true, err
		}

		if err = PublicationUpdateDateValidation(date); err != nil {
			return true, err
		}

		// todo переделать с storeData.ChannelID на storeData.PublicationID
		if err = b.publicationService.UpdateDeleteDate(ctx, storeData.ChannelID, date); err != nil {
			b.log.Error("isStoreExist::store.PublicationDeleteDateUpdate: %v", err)
		}

		// удаление происходит в [scheduled.go] в случае успешной отправки сообщения

	case store.PublicationSentDateUpdate:
		var (
			date time.Time
		)
		date, err = time.Parse(dto.Layout, update.Message.Text+":00 +0300")
		if err != nil {
			b.log.Error("isStoreExist::store.PublicationSentDateUpdate: %v", err)
			return true, err
		}

		if err = PublicationUpdateDateValidation(date); err != nil {
			return true, err
		}

		// todo переделать с storeData.ChannelID на storeData.PublicationID
		if err = b.publicationService.UpdatePublicationDate(ctx, storeData.ChannelID, date); err != nil {
			b.log.Error("isStoreExist::store.PublicationSentDateUpdate: %v", err)
		}
		if err == nil {
			b.log.Info("set publication date: date=%v channelID=%d", date, storeData.ChannelID)
			b.publicationArray.AppendPub(&store.PubData{
				PubDate:       date,
				PublicationID: storeData.ChannelID,
			})
		}

	default:
		return false, nil
	}

	if err == nil {
		b.response(storeData.OperationType, storeData.CurrentMsgID, storeData.PreferMsgID, storeData.ChannelID, update)
	}
	return true, err
}

func ConvertToMarkdownV2(text string, messageEntities []tgbotapi.MessageEntity) string {
	insertions := make(map[int]string)
	for _, e := range messageEntities {
		var before, after string
		if e.IsBold() {
			before = "*"
			after = "*"
		} else if e.IsItalic() {
			before = "_"
			after = "_"
		} else if e.Type == "underline" {
			before = "__"
			after = "__"
		} else if e.Type == "strikethrough" {
			before = "~"
			after = "~"
			//} else if e.Type == "spoiler" {
			//	before = "||"
			//	after = "||"
		} else if e.IsCode() {
			before = "`"
			after = "`"
		} else if e.IsPre() {
			before = "```" + e.Language
			after = "```"
		} else if e.IsTextLink() {
			before = "["
			after = "](" + e.URL + ")"
		}
		if before != "" {
			insertions[e.Offset] += before
			insertions[e.Offset+e.Length] += after
		}
	}

	input := []rune(text)
	var output []rune
	utf16pos := 0
	for _, c := range input {
		output = append(output, []rune(insertions[utf16pos])...)
		if _, has := needEscape[c]; has {
			output = append(output, '\\')
		}
		output = append(output, c)
		utf16pos += len(utf16.Encode([]rune{c}))
	}
	output = append(output, []rune(insertions[utf16pos])...)
	return string(output)
}
