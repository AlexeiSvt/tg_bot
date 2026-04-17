package handlers

import (
	"database/sql"
    "log"
    "tgbot/internal/states"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCallback(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, update tgbotapi.Update) {
	if update.CallbackQuery == nil {
		return
	}

	data := update.CallbackQuery.Data
	user := update.CallbackQuery.From
	chatID := update.CallbackQuery.Message.Chat.ID

	switch data {

	case CBDiscBS:
		handleDisciplineRules(bot, mgr, user.ID, chatID, "Brawl Stars", "bs", rulesBS)

	case CBDiscCR:
		handleDisciplineRules(bot, mgr, user.ID, chatID, "Clash Royale", "cr", rulesCR)

	case CBDiscCH:
		handleDisciplineRules(bot, mgr, user.ID, chatID, "Chess", "ch", rulesCH)

	case CBDiscTri:
		handleTriathlonStart(bot, mgr, user.ID, chatID)

	case CBTriBS:
		handleTriathlonGameSelect(bot, mgr, user.ID, chatID, "Brawl Stars")

	case CBTriCR:
		handleTriathlonGameSelect(bot, mgr, user.ID, chatID, "Clash Royale")

	case CBTriCH:
		handleTriathlonGameSelect(bot, mgr, user.ID, chatID, "Chess")

	case CBTriCheck:
		handleTriathlonCheck(bot, mgr, user.ID, chatID)

	case CBTriDone:
		handleTriathlonComplete(bot, mgr, user.ID, chatID)

	case CBMoreYes:
		handleMoreDisciplines(bot, mgr, user.ID, chatID)

	case CBMoreNo:
		handleRegistrationComplete(bot, mgr, user.ID, chatID)

	case CBTriConfirm, CBFinalConfirm:
		handleConfirmRegistration(bot, db, mgr, user.ID, chatID)

	case CBCancel:
		handleCancelRegistration(bot, mgr, user.ID, chatID)

	default:
		handleDefaultCallback(bot, mgr, user.ID, chatID, data)
	}

	bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

func handleDefaultCallback(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, data string) {
	if len(data) > len(CBPrefixOK) && data[:len(CBPrefixOK)] == CBPrefixOK {
		code := data[len(CBPrefixOK):]
		handleRulesOk(bot, mgr, userID, chatID, code)
		return
	}

	log.Printf("Unknown callback: %s from user %d", data, userID)
}
