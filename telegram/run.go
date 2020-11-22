package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go_telegram_notificator/db"
	telegram "go_telegram_notificator/telegram/utils"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
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
		if update.CallbackQuery != nil {
			if !checkAccess(update.CallbackQuery.Message.Chat.ID) {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Access denied!")
				bot.Send(msg)
				return
			}

			data := strings.Split(update.CallbackQuery.Data, ",")
			fmt.Println(data)
			fmt.Println(data[0])
			fmt.Println(data[1])
			if data[0] == "remove_user" {
				userId, _ := strconv.Atoi(data[1])
				chatId, _ := strconv.Atoi(update.CallbackQuery.ID)
				if db.RemoveUser(userId) {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "User removed!")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "User undefined!")
					bot.Send(msg)
				}
				//TODO remove message
				deleteConfig := tgbotapi.DeleteMessageConfig{ChatID: update.CallbackQuery.Message.Chat.ID, MessageID: chatId}
				bot.DeleteMessage(deleteConfig)
			}
		}
		if update.Message == nil {
			continue
		}

		//Проверяем что от пользователья пришло именно текстовое сообщение
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":
				fmt.Println(update.Message.Command())
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a Notification Bot")
				bot.Send(msg)
				fmt.Println(update)

			case "/allusers":
				//Присваиваем количество пользоватьелей использовавших бота в num переменную
				users := db.GetAllUsers()
				if err != nil {
					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
					bot.Send(msg)
				}
				var listUsers []string
				for _, user := range users {
					listUsers = append(listUsers, fmt.Sprintf("%s %s - %d", user.FirstName, user.LastName, user.ChatId))
				}
				//Создаем строку которая содержит колличество пользователей использовавших бота
				ans := fmt.Sprintf("People used Notification bot:\n%s ", telegram.ListAll(listUsers))

				//Отправлем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
				bot.Send(msg)

			case "/notify":
				if !checkAccess(update.Message.Chat.ID) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied!")
					bot.Send(msg)
					return
				}

				//Присваиваем количество пользоватьелей использовавших бота в num переменную
				users := db.GetAllUsers()
				if err != nil {
					//Отправлем сообщение
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
					bot.Send(msg)
				}
				for _, user := range users {
					PushMessage("It's test message", int64(user.ChatId))
				}
			case "/removeuser":
				if !checkAccess(update.Message.Chat.ID) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied!")
					bot.Send(msg)
					return
				}
				users := db.GetAllUsers()
				var numericKeyboard []tgbotapi.InlineKeyboardButton
				for _, user := range users {
					numericKeyboard = append(numericKeyboard, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintln(user.FirstName, user.LastName), fmt.Sprintf("remove_user,%d", user.ChatId)))
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select user:")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(numericKeyboard...),
				)
				bot.Send(msg)
			default:

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please use command!")
				if _, err := bot.Send(msg); err != nil {
					fmt.Println(err)
				}

			}
		} else if update.Message.Contact != nil {
			if !checkAccess(update.Message.Chat.ID) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied!")
				bot.Send(msg)
				return
			}

			if db.CreateUser(update.Message.Contact.FirstName, update.Message.Contact.LastName, update.Message.Contact.UserID) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("User %s %s created", update.Message.Contact.FirstName, update.Message.Contact.LastName))
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "User already created!")
				bot.Send(msg)
			}
		} else {

			//Отправлем сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}

}

func checkAccess(id int64) bool {
	chatId, _ := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	return chatId == id
}
