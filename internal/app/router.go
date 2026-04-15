package app

import (
	"database/sql"
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
	case "start":
		handlers.HandleStart(bot, mgr, update)
	case "help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используйте /start для регистрации, /cancel для отмены, /mystats для просмотра данных.")
		bot.Send(msg)
	case "cancel":
		mgr.Reset(update.Message.From.ID)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Регистрация отменена."))
	case "backup":
		if update.Message.Chat.ID == constants.AdminChatID {
			go backup.PerformBackup(bot, db)
			bot.Send(tgbotapi.NewMessage(constants.AdminChatID, "⏳ Создаю бэкап..."))
		}
	default:
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
	}
}