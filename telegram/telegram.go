package telegram

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const sendMsgMethod = "sendMessage"
const defaultTelegramUrl = "https://api.telegram.org"

var baseUrl string

type (
	Telegram struct {
		channel string
	}
)

func New(token, channel string) *Telegram {
	baseUrl = fmt.Sprintf("%s/bot%s/", defaultTelegramUrl, token)
	return &Telegram{channel: "@" + channel}
}

func (t *Telegram) SendMsg(text string) error {
	url := baseUrl + sendMsgMethod
	d := map[string]interface{}{
		"chat_id":                  t.channel,
		"text":                     text,
		"disable_web_page_preview": true,
	}
	_, err := resty.New().R().SetBody(d).Post(url)
	if err != nil {
		return err
	}
	return nil
}
