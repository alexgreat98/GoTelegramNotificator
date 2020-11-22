package main

import (
	"context"
	"go_telegram_notificator/db"
	"go_telegram_notificator/telegram"
)

func main() {
	ctx := context.Background()

	if err := db.Boot(&ctx); err != nil {
		panic(err)
	}
	telegram.Run()
}
