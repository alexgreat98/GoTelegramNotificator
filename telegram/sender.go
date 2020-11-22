package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
)

var Messages = make(chan Message)

func PushMessage(message string) {
	go func() { Messages <- Message{message} }()
}

func Sender(bot *tgbotapi.BotAPI) {
	for {
		select {
		case message := <-Messages:
			fmt.Println("text: ", message.Text())
			chatId, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
			msg := tgbotapi.NewMessage(chatId, message.Text())
			if _, err := bot.Send(msg); err != nil {
				fmt.Println(err)
			}
		}
	}
}
