package backup

import (
	"log"
	"os"
	"strconv"
	"tgbot/internal/constants"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendBackupFile(bot *tgbotapi.BotAPI, filename string) error {
	chatID, err := strconv.ParseInt(constants.AdminChatID, 10, 64)
	if err != nil {
		log.Fatalf("%s: %v", InvalidAdminID, err)
	}
	file := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filename))

	fileInfo, _ := os.Stat(filename)
	fileSize := float64(fileInfo.Size()) / 1024

	file.Caption = buildCaption(filename, fileSize)

	_, err = bot.Send(file)
	return err
}
