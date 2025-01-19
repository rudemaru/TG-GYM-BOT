package main

import (
	"fmt"
	"log"

	"github.com/rudemaru/TG-GYM-BOT/internal/bot"
	"github.com/rudemaru/TG-GYM-BOT/internal/config"
	"github.com/rudemaru/TG-GYM-BOT/internal/db"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	db.InitDB(cfg)

	bot, err := bot.NewBot(cfg.TelegramBotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Ready to serve o7")

	bot.Start()
}
