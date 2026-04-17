package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tgbot/internal/database"
	"tgbot/internal/models"
	"tgbot/internal/states"
	"tgbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleDisciplineRules(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, gameName, code, rules string) {
	bot.Send(tgbotapi.NewMessage(chatID, rules))

	m := tgbotapi.NewMessage(chatID, RulesConfirmPrompt)
	m.ReplyMarkup = utils.RulesOkButton(code)
	bot.Send(m)

	mgr.SetState(userID, states.ReadingRules)
}

func handleTriathlonStart(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
	s := mgr.Get(userID)

	s.TriGames = map[string]bool{
		"Brawl Stars":  true,
		"Clash Royale": true,
		"Chess":        true,
	}

	msg := tgbotapi.NewMessage(chatID, TriathlonRules)
	msg.ReplyMarkup = getTriathlonKeyboard(s.Temp.Disciplines)
	bot.Send(msg)

	mgr.SetState(userID, states.TriathlonSelect)
}

func handleTriathlonGameSelect(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, gameName string) {
	s := mgr.Get(userID)
	s.CurrentGame = gameName

	mgr.SetState(userID, states.EnteringNick)

	bot.Send(tgbotapi.NewMessage(chatID,
		fmt.Sprintf(EnterNick, gameName)))
}

func handleTriathlonCheck(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)
    msg := tgbotapi.NewMessage(chatID, getTriathlonStatus(s.Temp.Disciplines))
    msg.ReplyMarkup = getTriathlonKeyboard(s.Temp.Disciplines)
    bot.Send(msg)
}

func handleTriathlonComplete(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
	s := mgr.Get(userID)

	if !isTriathlonComplete(s.Temp.Disciplines) {
		bot.Send(tgbotapi.NewMessage(chatID, TriathlonIncomplete))
		return
	}

	showConfirmationPreview(bot, userID, chatID, s.Temp, "tri_confirm")
}

func handleMoreDisciplines(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
	s := mgr.Get(userID)

	s.TriGames = nil

	mgr.SetState(userID, states.ChoosingDiscipline)

	msg := tgbotapi.NewMessage(chatID, ChooseNextGame)
	msg.ReplyMarkup = utils.DisciplineKeyboard()
	bot.Send(msg)
}

func handleRegistrationComplete(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)

    showConfirmationPreview(bot, userID, chatID, s.Temp, "final_confirm")
}

func handleRulesOk(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, code string) {
	s := mgr.Get(userID)

	gameMap := map[string]string{
		"bs": "Brawl Stars",
		"cr": "Clash Royale",
		"ch": "Chess",
	}

	gameName, ok := gameMap[code]
	if !ok {
		log.Printf("Unknown game code: %s", code)
		return
	}

	s.CurrentGame = gameName
	mgr.SetState(userID, states.EnteringNick)

	bot.Send(tgbotapi.NewMessage(chatID,
		fmt.Sprintf(EnterNick, gameName)))
}


func showConfirmationPreview(bot *tgbotapi.BotAPI, userID, chatID int64, u *models.User, confirmCode string) {
	var preview strings.Builder; preview.WriteString(ConfirmationHeader +
		fmt.Sprintf(ConfirmationBody,
			u.FirstName, u.LastName, u.Class))

	for game, gd := range u.Disciplines {
		if game == "Chess" {
			preview .WriteString(fmt.Sprintf("  🔸 %s: %s\n", game, gd.Nick))
		} else {
			preview .WriteString(fmt.Sprintf("  🔸 %s: %s | %s\n", game, gd.Nick, gd.Tag))
		}
	}

	preview .WriteString(ConfirmationFooter)

	msg := tgbotapi.NewMessage(chatID, preview.String())
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ConfirmButton, confirmCode),
			tgbotapi.NewInlineKeyboardButtonData(CancelButton, "cancel_reg"),
		),
	)

	bot.Send(msg)
}

func handleConfirmRegistration(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, userID, chatID int64) {
	s := mgr.Get(userID)
	s.Temp.TelegramID = userID

	if err := database.SaveUser(db, s.Temp); err != nil {
		log.Printf("Error saving user: %v", err)
		bot.Send(tgbotapi.NewMessage(chatID, SaveError))
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, formatSummary(s.Temp)))
	mgr.Reset(userID)
}

func handleCancelRegistration(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
	mgr.Reset(userID)

	msg := tgbotapi.NewMessage(chatID, RegistrationCancelled)
	bot.Send(msg)
}