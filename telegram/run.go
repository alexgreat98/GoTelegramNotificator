package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go_telegram_notificator/db"
	"log"
	"os"
	"reflect"
)

func Run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("telegram_error: ", err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go Sender(bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//Проверяем что от пользователья пришло именно текстовое сообщение
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				//Отправлем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a Notification Bot")
				if err := db.CollectData(update.Message.Chat.UserName, update.Message.Chat.ID, update.Message.Text, []string{msg.Text}); err != nil {
					fmt.Println(err)
				}
				fmt.Println(msg)
				fmt.Println(update)
				bot.Send(msg)

			case "/allusers":
				if os.Getenv("DB_SWITCH") == "on" {

					//Присваиваем количество пользоватьелей использовавших бота в num переменную
					users, err := db.GetAllUsers()
					if err != nil {
						//Отправлем сообщение
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						bot.Send(msg)
					}

					//Создаем строку которая содержит колличество пользователей использовавших бота
					ans := fmt.Sprintf("%s peoples used Notification bot", users)

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
					bot.Send(msg)
				} else {

					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database not connected, so i can't say you how many peoples used me.")
					bot.Send(msg)
				}
			default:

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please use command!")
				fmt.Println(update.Message.Text)
				fmt.Println(update.Message.Contact)
				if _, err := bot.Send(msg); err != nil {
					fmt.Println(err)
				}

			}
		} else if update.Message.Contact != nil {
			fmt.Println(update.Message.Text)
			fmt.Println(update.Message.Contact.UserID)
		} else {

			//Отправлем сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}
