package backup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"tgbot/internal/constants"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PerformBackup(bot *tgbotapi.BotAPI, db *sql.DB) {
	filename := generateFilename()

	chatID, err := strconv.ParseInt(constants.AdminChatID, 10, 64)
	if err != nil {
		log.Fatalf("%s: %v",InvalidAdminID, err)
	}

	err = ExportToCSV(db, filename)
	if err != nil {
		log.Printf("%s, %v", BackupCreationMistake, err)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %v", BackupCreationMistake, err))
		bot.Send(msg)
		return
	}
	defer os.Remove(filename)

	err = SendBackupFile(bot, filename)
	if err != nil {
		log.Printf("%s, %v", BackupSendingMistake, err)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("%s: %v", BackupSendingMistake, err))
		bot.Send(msg)
		return
	}

	log.Printf("%s, %s", BackupWasSentSuccessfully, filename)
}

func generateFilename() string {
	return fmt.Sprintf("backup_etriathlon_%s.csv",
		time.Now().Format("2006-01-02_15-04-05"))
}
