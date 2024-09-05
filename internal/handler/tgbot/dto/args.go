package dto

import (
	"encoding/json"
	"time"
)

type Button struct {
	ButtonUrl  *string `json:"текст_кнопки"`
	ButtonText *string `json:"ссылка_кнопки"`
}

type PublicationCreate struct {
	PublicationDate time.Time  `json:"дата_публикации"`
	DeleteDate      *time.Time `json:"дата_удаления"`
	Button          Button     `json:"кнопка"`
}

const Layout = "2006-01-02 15:04:05 -0700"

func (c *PublicationCreate) UnmarshalJSON(b []byte) (err error) {
	var jsonMap map[string]interface{}
	err = json.Unmarshal(b, &jsonMap)
	if err != nil {
		return
	}

	publicationDateStr, _ := jsonMap["дата_публикации"].(string)
	c.PublicationDate, err = time.Parse(Layout, publicationDateStr+":00 +0300")
	if err != nil {
		return
	}

	deleteDateStr, ok := jsonMap["дата_удаления"].(string)
	if ok {
		c.DeleteDate = new(time.Time)
		*c.DeleteDate, err = time.Parse(Layout, deleteDateStr+":00 +0300")
		if err != nil {
			return
		}
	}

	button, ok := jsonMap["кнопка"].(map[string]interface{})
	if ok {
		if btnUrl, urlExist := button["текст_кнопки"].(string); urlExist {
			c.Button.ButtonUrl = new(string)
			*c.Button.ButtonUrl = btnUrl
		}
		if btnText, textExist := button["ссылка_кнопки"].(string); textExist {
			c.Button.ButtonText = new(string)
			*c.Button.ButtonText = btnText
		}
	}

	return
}
