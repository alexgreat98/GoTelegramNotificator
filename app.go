package main

import (
	"go_telegram_notificator/db"
	"go_telegram_notificator/telegram"
)

func main() {

	if _, err := db.Boot(); err != nil {
		panic(err)
	}
	telegram.Run()
}
