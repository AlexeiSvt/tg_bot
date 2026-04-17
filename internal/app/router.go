package app

import (
	"database/sql"
	"log"
	"strconv"
	"tgbot/internal/backup"
	"tgbot/internal/constants"
	"tgbot/internal/handlers"
	"tgbot/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleUpdate(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, update tgbotapi.Update) {
	if update.Message != nil {
		HandleMessage(bot, db, mgr, update)
		return
	}

	if update.CallbackQuery != nil {
		handlers.HandleCallback(bot, db, mgr, update)
	}
}

func HandleMessage(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, update tgbotapi.Update) {
	if !update.Message.IsCommand() {
		handlers.HandleMessage(bot, db, mgr, update)
		return
	}

	switch update.Message.Command() {
	case Start:
		handlers.HandleStart(bot, mgr, update)
	case Help:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, StartMessage)
		bot.Send(msg)
	case Cancel:
		mgr.Reset(update.Message.From.ID)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, CancelRegistration))
	case Backup:
		chatID, err := strconv.ParseInt(constants.AdminChatID, 10, 64)
		if err != nil {
			log.Fatalf("%s: %v", InvalidChatID, err)
		}

		if update.Message.Chat.ID == chatID {
			go backup.PerformBackup(bot, db)
			bot.Send(tgbotapi.NewMessage(chatID, CreateBackup))
		}
	default:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, UnknownCommand))
	}
}
