package backup

import (
	"database/sql"
	"log"
	"strconv"
	"tgbot/internal/constants"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBackupRoutine(bot *tgbotapi.BotAPI, db *sql.DB) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	chatID, err := strconv.ParseInt(constants.AdminChatID, 10, 64)
	if err != nil {
		log.Fatalf("%s: %v", InvalidAdminID, err)
	}

	msg := tgbotapi.NewMessage(chatID, StartSystemBackup)
	bot.Send(msg)

	for range ticker.C {
		PerformBackup(bot, db)
	}
}
