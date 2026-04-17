package handlers

import (
	"database/sql"
	"fmt"
	"tgbot/internal/states"
	"tgbot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleNameInput(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, text string) {
	s := mgr.Get(userID)
	s.Temp.FirstName = text

	mgr.SetState(userID, states.WaitingLastName)
	bot.Send(tgbotapi.NewMessage(chatID, MsgAskLastName))
}

func handleLastNameInput(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, text string) {
	s := mgr.Get(userID)
	s.Temp.LastName = text

	mgr.SetState(userID, states.WaitingClass)
	bot.Send(tgbotapi.NewMessage(chatID, MsgAskClass))
}

func handleClassInput(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, text string) {
	s := mgr.Get(userID)
	s.Temp.Class = text

	mgr.SetState(userID, states.ChoosingDiscipline)

	msg := tgbotapi.NewMessage(chatID, MsgChooseDisc)
	msg.ReplyMarkup = utils.DisciplineKeyboard()
	bot.Send(msg)
}

func handleNickInput(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, text string) {
	s := mgr.Get(userID)

	if s.CurrentGame == "" {
		bot.Send(tgbotapi.NewMessage(chatID, MsgNoGame))
		mgr.SetState(userID, states.ChoosingDiscipline)
		return
	}

	gd := s.Temp.Disciplines[s.CurrentGame]
	gd.Nick = text
	s.Temp.Disciplines[s.CurrentGame] = gd

	if s.CurrentGame == MsgChessDone {
		handlePostChessNick(bot, mgr, userID, chatID)
		return
	}

	mgr.SetState(userID, states.EnteringTag)
	bot.Send(tgbotapi.NewMessage(chatID,
		fmt.Sprintf(MsgEnterTag, s.CurrentGame)))
}

func handlePostChessNick(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
	s := mgr.Get(userID)

	isTriathlon := len(s.TriGames) > 0

	if isTriathlon {
		msg := tgbotapi.NewMessage(chatID, MsgChessSaved)
		msg.ReplyMarkup = getTriathlonKeyboard(s.Temp.Disciplines)
		bot.Send(msg)

		mgr.SetState(userID, states.TriathlonSelect)
		return
	}

	askMoreDisciplines(bot, mgr, userID, chatID)
}

func handleTagInput(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, userID, chatID int64, text string) {
	s := mgr.Get(userID)

	if !utils.ValidateTag(text) {
		bot.Send(tgbotapi.NewMessage(chatID, MsgInvalidTag))
		return
	}

	gd := s.Temp.Disciplines[s.CurrentGame]
	gd.Tag = text
	s.Temp.Disciplines[s.CurrentGame] = gd

	isTriathlon := len(s.TriGames) > 0

	if isTriathlon {
		msg := tgbotapi.NewMessage(chatID,
			fmt.Sprintf("✅ Данные для %s сохранены!\n\nВыберите следующую игру:", s.CurrentGame))
		msg.ReplyMarkup = getTriathlonKeyboard(s.Temp.Disciplines)
		bot.Send(msg)

		mgr.SetState(userID, states.TriathlonSelect)
		return
	}

	askMoreDisciplines(bot, mgr, userID, chatID)
}
