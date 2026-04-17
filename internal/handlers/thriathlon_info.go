package handlers

import (
	"fmt"
	"strings"
	"tgbot/internal/models"
	"tgbot/internal/states"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func askMoreDisciplines(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "more_yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет, завершить", "more_no"),
		),
	)

	msg := tgbotapi.NewMessage(chatID, MsgMoreGames)
	msg.ReplyMarkup = kb

	bot.Send(msg)
	mgr.SetState(userID, states.ChoosingDiscipline)
}

func getTriathlonKeyboard(disciplines map[string]models.GameData) tgbotapi.InlineKeyboardMarkup {
	games := []struct {
		name string
		code string
	}{
		{"Brawl Stars", "tri_bs"},
		{"Clash Royale", "tri_cr"},
		{"Chess", "tri_ch"},
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, game := range games {
		status := "⬜"
		if gd, ok := disciplines[game.name]; ok && gd.Nick != "" {
			status = "✅"
		}

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s %s", status, game.name),
					game.code,
				),
			),
		)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Статус", "tri_check"),
		),
	)

	if isTriathlonComplete(disciplines) {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Завершить регистрацию", "tri_done"),
			),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

type TriGame struct {
	Name    string
	NeedTag bool
}

var triGames = []TriGame{
	{"Brawl Stars", true},
	{"Clash Royale", true},
	{"Chess", false},
}

func isTriathlonComplete(disciplines map[string]models.GameData) bool {
	for _, g := range triGames {
		gd, ok := disciplines[g.Name]
		if !ok || gd.Nick == "" {
			return false
		}

		if g.NeedTag && gd.Tag == "" {
			return false
		}
	}
	return true
}

func getTriathlonStatus(disciplines map[string]models.GameData) string {
	var status strings.Builder

	status.WriteString(TriStatusHeader)

	for _, g := range triGames {
		gd, ok := disciplines[g.Name]

		icon := "⬜"
		details := TriNotFilled

		if ok && gd.Nick != "" {
			icon = "✅"

			if g.NeedTag {
				details = fmt.Sprintf(TriNickTag, gd.Nick, gd.Tag)
			} else {
				details = fmt.Sprintf(TriNickOnly, gd.Nick)
			}
		}

		status.WriteString(fmt.Sprintf(TriStatusLine, icon, g.Name, details))
	}

	return status.String()
}
