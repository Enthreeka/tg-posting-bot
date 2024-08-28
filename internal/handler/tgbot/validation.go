package tgbot

import (
	"errors"
	"github.com/Enthreeka/tg-posting-bot/internal/handler/tgbot/dto"
	"time"
)

// todo вынести ошибки в custom err
func PublicationCreateValidation(msg dto.PublicationCreate) error {
	if msg.Text == "" {
		return errors.New("текст должен быть заполнен")
	}

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	// Преобразуем все даты в часовой пояс "Europe/Moscow"
	publicationDate := msg.PublicationDate.In(loc)
	deleteDate := msg.DeleteDate
	if deleteDate != nil {
		*deleteDate = deleteDate.In(loc)
	}

	// Проверяем, является ли время публикации раньше текущего времени
	if publicationDate.Before(time.Now().In(loc)) {
		return errors.New("время публикации раньше чем текущее время по Europe/Moscow")
	}

	// Проверяем, является ли время удаления раньше текущего времени
	if deleteDate != nil && deleteDate.Before(time.Now().In(loc)) {
		return errors.New("время удаления раньше чем текущее время по Europe/Moscow")
	}

	// Проверяем, является ли время удаления раньше времени публикации
	if deleteDate != nil && deleteDate.Before(publicationDate) {
		return errors.New("время удаления раньше чем время публикации по Europe/Moscow")
	}

	if msg.Button != (dto.Button{}) {
		if msg.Button.ButtonUrl == nil {
			return errors.New("отсутствует текст для кнопки")
		}
		if msg.Button.ButtonText == nil {
			return errors.New("отсутствует ссылка для кнопки")
		}
	}

	return nil
}

func PublicationUpdateDateValidation(date time.Time) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	if date.Before(time.Now().In(loc)) {
		return errors.New("время публикации раньше чем текущее время по Europe/Moscow")
	}

	return nil
}
