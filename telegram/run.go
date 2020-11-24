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
	var notificationChatId int64 = 0
	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	go Sender(bot)

	for update := range updates {
		if update.CallbackQuery != nil {
			if !checkAccess(update.CallbackQuery.From.ID) {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Access denied!")
				bot.Send(msg)
				return
			}

			data := strings.Split(update.CallbackQuery.Data, ",")

			if data[0] == "remove_user" {
				userId, _ := strconv.Atoi(data[1])
				chatId := update.CallbackQuery.Message.MessageID
				if db.RemoveUser(userId) {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "User removed!")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "User undefined!")
					bot.Send(msg)
				}
				deleteConfig := tgbotapi.DeleteMessageConfig{ChatID: update.CallbackQuery.Message.Chat.ID, MessageID: chatId}
				_, err := bot.DeleteMessage(deleteConfig)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if update.Message == nil {
			continue
		}

		//Проверяем что от пользователья пришло именно текстовое сообщение
		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":
				fmt.Println(update.Message.From.UserName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a Notification Bot")
				bot.Send(msg)

			case "/register":
				var username string
				if update.Message.From.UserName == "" {
					username = update.Message.From.FirstName
				}
				if db.CreateUser(update.Message.From.FirstName, update.Message.From.LastName, update.Message.From.ID, username) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("User %s %s created", update.Message.From.FirstName, update.Message.From.LastName))
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("User %s already created!", update.Message.From.FirstName))
					bot.Send(msg)
				}
			case "/allusers":
				//Присваиваем количество пользоватьелей использовавших бота в num переменную
				users := db.GetAllUsers()

				var ans string

				if len(users) <= 0 {
					ans = "User list is empty."
				} else {
					var listUsers []string
					for _, user := range users {
						listUsers = append(listUsers, fmt.Sprintf("%s %s - %d", user.FirstName, user.LastName, user.ChatId))
					}
					ans = fmt.Sprintf("People used Notification bot:\n%s ", telegram.ListAll(listUsers))
				}

				//Отправлем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, ans)
				bot.Send(msg)

			case "/notify":
				if !checkAccess(update.Message.From.ID) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied!")
					bot.Send(msg)
					return
				}

				if notificationChatId == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Chat undefined!")
					bot.Send(msg)
					return
				}
				//Присваиваем количество пользоватьелей использовавших бота в num переменную
				users := db.GetAllUsers()
				for _, user := range users {
					PushMessage(fmt.Sprintf("[%s %s](tg://user?id=%d) Hey man", user.FirstName, user.LastName, user.ChatId), notificationChatId)
				}
			case "/setchat":
				if !checkAccess(update.Message.From.ID) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied!")
					bot.Send(msg)
					return
				}

				notificationChatId = update.Message.Chat.ID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This chat is used for notification.")
				bot.Send(msg)
			case "/removeuser":
				if !checkAccess(update.Message.From.ID) {
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
		} else {

			//Отправлем сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please use command!")
			bot.Send(msg)
		}
	}

}

func checkAccess(id int) bool {
	chatId, _ := strconv.Atoi(os.Getenv("ADMIN_ID"))
	return chatId == id
}
