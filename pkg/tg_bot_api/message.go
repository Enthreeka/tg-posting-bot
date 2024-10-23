package tg_bot_api

import (
	"github.com/Enthreeka/tg-posting-bot/internal/entity"
	"github.com/Enthreeka/tg-posting-bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"strings"
)

type Message interface {
	SendNewMessage(chatID int64, markup *tgbotapi.InlineKeyboardMarkup, text string) (int, error)
	SendEditMessage(chatID int64, messageID int, markup *tgbotapi.InlineKeyboardMarkup, text string) (int, error)
	SendDocument(chatID int64, fileName string, fileIDBytes *[]byte, text string) (int, error)
	SendMessageToUser(chatID int64, publication *entity.Publication) (int, error)
	SendMessageToChannel(username string, publication *entity.Publication) error
	DeleteMessage(chatID int64, messageID int) error
}

type TelegramMsg struct {
	log *logger.Logger
	bot *tgbotapi.BotAPI
}

func NewMessageSetting(bot *tgbotapi.BotAPI, log *logger.Logger) *TelegramMsg {
	return &TelegramMsg{
		bot: bot,
		log: log,
	}
}

func (t *TelegramMsg) SendNewMessage(chatID int64, markup *tgbotapi.InlineKeyboardMarkup, text string) (int, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true
	if markup != nil {
		msg.ReplyMarkup = &markup
	}

	sendMsg, err := t.bot.Send(msg)
	if err != nil {
		t.log.Error("failed to send message", zap.Error(err))
		return 0, err
	}

	return sendMsg.MessageID, nil
}

func (t *TelegramMsg) SendEditMessage(chatID int64, messageID int, markup *tgbotapi.InlineKeyboardMarkup, text string) (int, error) {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.DisableWebPagePreview = true

	if markup != nil {
		msg.ReplyMarkup = markup
	}

	sendMsg, err := t.bot.Send(msg)
	if err != nil {
		t.log.Error("failed to send msg: %v", err)
		return 0, err
	}

	return sendMsg.MessageID, nil
}

func (t *TelegramMsg) SendDocument(chatID int64, fileName string, fileIDBytes *[]byte, text string) (int, error) {
	msg := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: *fileIDBytes,
	})
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	msg.Caption = text

	sendMsg, err := t.bot.Send(msg)
	if err != nil {
		t.log.Error("failed to send msg: %v", err)
		return 0, err
	}

	return sendMsg.MessageID, nil
}

func (t *TelegramMsg) SendMessageToUser(chatID int64, publication *entity.Publication) (int, error) {
	if publication.Image != nil {
		publicationPhotoPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(*publication.Image))
		msg := tgbotapi.NewPhoto(chatID, publicationPhotoPhoto.Media)
		buttonMarkup := buttonQualifier(publication.ButtonUrl, publication.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if publication.Text != "" {
			msg.Caption = publication.Text
		}
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		sendMsg, err := t.bot.Send(msg)
		if err != nil {
			t.log.Error("failed to send message: %v", err)
			return 0, err
		}
		return sendMsg.MessageID, nil
	}

	msg := tgbotapi.NewMessage(chatID, "")
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	buttonMarkup := buttonQualifier(publication.ButtonUrl, publication.ButtonText)
	if buttonMarkup != nil {
		msg.ReplyMarkup = &buttonMarkup
	}
	if publication.Text != "" {
		msg.Text = publication.Text
	}

	sendMsg, err := t.bot.Send(msg)
	if err != nil {
		t.log.Error("failed to send message", err)
		return 0, err
	}

	return sendMsg.MessageID, nil
}

func (t *TelegramMsg) DeleteMessage(chatID int64, messageID int) error {
	resp, err := t.bot.Request(tgbotapi.NewDeleteMessage(chatID, messageID))
	if nil != err || !resp.Ok {
		t.log.Error("failed to delete message id %d (%s): %v", messageID, string(resp.Result), err)
	}
	return err
}

func (t *TelegramMsg) SendMessageToChannel(username string, publication *entity.Publication) error {
	if publication.Image != nil {
		publicationPhoto := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(*publication.Image))
		msg := tgbotapi.NewPhotoToChannel(username, publicationPhoto.Media)
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		buttonMarkup := buttonQualifier(publication.ButtonUrl, publication.ButtonText)
		if buttonMarkup != nil {
			msg.ReplyMarkup = &buttonMarkup
		}
		if publication.Text != "" {
			msg.Caption = publication.Text
		}

		if _, err := t.bot.Send(msg); err != nil {
			t.log.Error("failed to send message: %v", err)
			return err
		}
		return nil
	}

	msg := tgbotapi.NewMessageToChannel(username, "")
	msg.DisableWebPagePreview = true
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	buttonMarkup := buttonQualifier(publication.ButtonUrl, publication.ButtonText)
	if buttonMarkup != nil {
		msg.ReplyMarkup = &buttonMarkup
	}
	if publication.Text != "" {
		msg.Text = publication.Text
	}

	if _, err := t.bot.Send(msg); err != nil {
		t.log.Error("failed to send message", err)
		return err
	}

	return nil
}

func buttonQualifier(buttonURL *string, buttonText *string) *tgbotapi.InlineKeyboardMarkup {
	if buttonURL != nil && buttonText != nil {
		var (
			btnText string
			btnURL  string
		)

		btnText = *buttonText
		btnURL = *buttonURL

		button := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(btnText, btnURL)),
		)
		return &button
	}
	return nil
}

func escapeSpecialCharacters(text string) string {
	specialChars := []string{"<", ">", "[", "]", "(", ")", "~", "`", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, `\`+char)
	}
	return text
}
