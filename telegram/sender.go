package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Messages = make(chan Message)

func PushMessage(message string, chatId int64) {
	go func() { Messages <- Message{message, chatId} }()
}

func Sender(bot *tgbotapi.BotAPI) {
	for {
		select {
		case message := <-Messages:
			fmt.Println("text: ", message.Text())
			msg := tgbotapi.NewMessage(message.ChatId(), message.Text())
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		}
	}
}
