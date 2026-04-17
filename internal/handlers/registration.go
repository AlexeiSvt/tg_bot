package handlers

import (
	"database/sql"
    "log"
    "tgbot/internal/states"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, update tgbotapi.Update) {
	if update.Message == nil || update.Message.From == nil {
		return
	}

	user := update.Message.From
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	s := mgr.Get(user.ID)

	switch s.State {
	case states.WaitingName:
		handleNameInput(bot, mgr, user.ID, chatID, text)

	case states.WaitingLastName:
		handleLastNameInput(bot, mgr, user.ID, chatID, text)

	case states.WaitingClass:
		handleClassInput(bot, mgr, user.ID, chatID, text)

	case states.EnteringNick:
		handleNickInput(bot, mgr, user.ID, chatID, text)

	case states.EnteringTag:
		handleTagInput(bot, db, mgr, user.ID, chatID, text)

	default:
		log.Printf("Unhandled state: %v for user %d", s.State, user.ID)
	}
}
