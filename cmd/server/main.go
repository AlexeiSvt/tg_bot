package main

import (
	"log"
	"tgbot/internal/app"
	"tgbot/internal/backup"
	"tgbot/internal/config"
	"tgbot/internal/database"
	"tgbot/internal/healthcheck"
	"tgbot/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("bot init: %v", err)
	}

	db, err := database.Open(cfg.DBDSN)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	mgr := states.NewManager()

	go backup.StartBackupRoutine(bot, db)
	go healtcheck.StartHealthCheckServer()

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 30
	updates := bot.GetUpdatesChan(ucfg)

	log.Println("Bot started successfully!")

	for update := range updates {
		app.HandleUpdate(bot, db, mgr, update)
	}
}
