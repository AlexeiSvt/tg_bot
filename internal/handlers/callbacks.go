package handlers

import (
    "database/sql"
    "fmt"
    "log"

    "tgbot/internal/database"
    "tgbot/internal/models"
    "tgbot/internal/states"
    "tgbot/internal/utils"

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
    // Обработка выбора одиночной дисциплины
    case "disc_bs":
        handleDisciplineRules(bot, mgr, user.ID, chatID, "Brawl Stars", "bs", rulesBS)
    case "disc_cr":
        handleDisciplineRules(bot, mgr, user.ID, chatID, "Clash Royale", "cr", rulesCR)
    case "disc_ch":
        handleDisciplineRules(bot, mgr, user.ID, chatID, "Chess", "ch", rulesCH)

    // Обработка триатлона
    case "disc_tri":
        handleTriathlonStart(bot, mgr, user.ID, chatID)

    // Обработка выбора игры в триатлоне
    case "tri_bs":
        handleTriathlonGameSelect(bot, mgr, user.ID, chatID, "Brawl Stars")
    case "tri_cr":
        handleTriathlonGameSelect(bot, mgr, user.ID, chatID, "Clash Royale")
    case "tri_ch":
        handleTriathlonGameSelect(bot, mgr, user.ID, chatID, "Chess")

    // Управление триатлоном
    case "tri_check":
        handleTriathlonCheck(bot, mgr, user.ID, chatID)
    case "tri_done":
        handleTriathlonComplete(bot, mgr, user.ID, chatID)

    // Управление обычной регистрацией
    case "more_yes":
        handleMoreDisciplines(bot, mgr, user.ID, chatID)
    case "more_no":
        handleRegistrationComplete(bot, mgr, user.ID, chatID)

    // Финальное подтверждение
    case "tri_confirm", "final_confirm":
        handleConfirmRegistration(bot, db, mgr, user.ID, chatID)

    // Отмена регистрации
    case "cancel_reg":
        handleCancelRegistration(bot, mgr, user.ID, chatID)

    // Обработка подтверждения правил (ok_bs, ok_cr, ok_ch)
    default:
        if len(data) > 3 && data[:3] == "ok_" {
            code := data[3:]
            handleRulesOk(bot, mgr, user.ID, chatID, code)
        } else {
            log.Printf("Unknown callback: %s from user %d", data, user.ID)
        }
    }

    // Ответ на callback запрос (убирает часики в Telegram)
    bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

// handleDisciplineRules показывает правила выбранной дисциплины
func handleDisciplineRules(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, gameName, code, rules string) {
    bot.Send(tgbotapi.NewMessage(chatID, rules))

    m := tgbotapi.NewMessage(chatID, "Нажмите кнопку ниже, если ознакомились с правилами:")
    m.ReplyMarkup = utils.RulesOkButton(code)
    bot.Send(m)

    mgr.SetState(userID, states.ReadingRules)
}

// handleTriathlonStart инициализирует регистрацию на триатлон
func handleTriathlonStart(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)
    // Отмечаем, что это триатлон — инициализируем TriGames
    s.TriGames = make(map[string]bool)
    s.TriGames["Brawl Stars"] = true
    s.TriGames["Clash Royale"] = true
    s.TriGames["Chess"] = true

    msg := tgbotapi.NewMessage(chatID,
        "🏆 ПРАВИЛА ТРИАТЛОНА\n\n"+
            "Вы участвуете во всех трёх играх:\n"+
            "• Brawl Stars\n"+
            "• Clash Royale\n"+
            "• Chess (Шахматы)\n\n"+
            "Для каждой игры необходимо ввести ник и тег игрока.\n\n"+
            "Выберите игру для ввода данных:")
    msg.ReplyMarkup = getTriathlonKeyboard(s.Temp.Disciplines)
    bot.Send(msg)

    mgr.SetState(userID, states.TriathlonSelect)
}

// handleTriathlonGameSelect переводит пользователя на ввод ника для выбранной игры
func handleTriathlonGameSelect(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64, gameName string) {
    s := mgr.Get(userID)
    s.CurrentGame = gameName
    mgr.SetState(userID, states.EnteringNick)

    bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Введите ваш ник в %s:", gameName)))
}

// handleTriathlonCheck показывает текущий статус заполнения
func handleTriathlonCheck(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)
    msg := tgbotapi.NewMessage(chatID, getTriathlonStatus(s.Temp.Disciplines))
    msg.ReplyMarkup = getTriathlonKeyboard(s.Temp.Disciplines)
    bot.Send(msg)
}

// handleTriathlonComplete завершает регистрацию на триатлон
func handleTriathlonComplete(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)
    if !isTriathlonComplete(s.Temp.Disciplines) {
        bot.Send(tgbotapi.NewMessage(chatID, "❌ Необходимо заполнить данные для всех трёх игр!"))
        return
    }

    // Показываем превью с просьбой подтвердить
    showConfirmationPreview(bot, userID, chatID, s.Temp, "tri_confirm")
}

// handleMoreDisciplines показывает оставшиеся дисциплины
func handleMoreDisciplines(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)
    // Очищаем флаг триатлона при выборе "Да"
    s.TriGames = nil

    mgr.SetState(userID, states.ChoosingDiscipline)
    msg := tgbotapi.NewMessage(chatID, "Выберите следующую игру:")
    msg.ReplyMarkup = utils.DisciplineKeyboard()
    bot.Send(msg)
}

// handleRegistrationComplete завершает регистрацию пользователя
func handleRegistrationComplete(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)

    // Показываем превью с просьбой подтвердить
    showConfirmationPreview(bot, userID, chatID, s.Temp, "final_confirm")
}

// handleRulesOk обрабатывает подтверждение правил
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
    bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Введите ваш ник в %s:", gameName)))
}

// showConfirmationPreview показывает превью данных и просит подтверждение
func showConfirmationPreview(bot *tgbotapi.BotAPI, userID, chatID int64, u *models.User, confirmCode string) {
    preview := fmt.Sprintf(
        "📋 ПРОВЕРКА ДАННЫХ\n\n"+
            "Пожалуйста, внимательно проверьте введённую информацию:\n\n"+
            "👤 Имя: %s\n"+
            "👤 Фамилия: %s\n"+
            "📚 Класс: %s\n\n"+
            "🎮 Дисциплины:\n",
        u.FirstName, u.LastName, u.Class)

    for game, gd := range u.Disciplines {
        if game == "Chess" {
            preview += fmt.Sprintf("  🔸 %s: %s\n", game, gd.Nick)
        } else {
            preview += fmt.Sprintf("  🔸 %s: %s | %s\n", game, gd.Nick, gd.Tag)
        }
    }

    preview += "\n" +
        "✅ Я подтверждаю правильность введённой информации и согласен с её обработкой в соответствии с правилами турнира eTriathlon 2026.\n\n" +
        "⚠️ Если вы обнаружили ошибку, удалите чат с ботом и создайте новый для корректировки данных."

    msg := tgbotapi.NewMessage(chatID, preview)
    msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("✅ Всё верно, завершить", confirmCode),
            tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", "cancel_reg"),
        ),
    )
    bot.Send(msg)
}

// handleConfirmRegistration обрабатывает финальное подтверждение
func handleConfirmRegistration(bot *tgbotapi.BotAPI, db *sql.DB, mgr *states.Manager, userID, chatID int64) {
    s := mgr.Get(userID)
    s.Temp.TelegramID = userID

    if err := database.SaveUser(db, s.Temp); err != nil {
        log.Printf("Error saving user: %v", err)
        bot.Send(tgbotapi.NewMessage(chatID, "❌ Ошибка при сохранении данных. Попробуйте позже."))
        return
    }

    bot.Send(tgbotapi.NewMessage(chatID, formatSummary(s.Temp)))
    mgr.Reset(userID)
}

// handleCancelRegistration отменяет регистрацию
func handleCancelRegistration(bot *tgbotapi.BotAPI, mgr *states.Manager, userID, chatID int64) {
    mgr.Reset(userID)
    msg := tgbotapi.NewMessage(chatID,
        "❌ Регистрация отменена.\n\n"+
            "Для начала новой регистрации введите /start")
    bot.Send(msg)
}

// formatSummary форматирует итоговое сообщение с данными пользователя
func formatSummary(u *models.User) string {
    summary := fmt.Sprintf("✅ Регистрация завершена!\n\n"+
        "Ваши данные:\n"+
        "📝 Имя: %s\n"+
        "📝 Фамилия: %s\n"+
        "📚 Класс: %s\n\n"+
        "🎮 Дисциплины:\n", u.FirstName, u.LastName, u.Class)

    for game, gd := range u.Disciplines {
        if game == "Chess" {
            summary += fmt.Sprintf("  🔸 %s: %s\n", game, gd.Nick)
        } else {
            summary += fmt.Sprintf("  🔸 %s: %s | %s\n", game, gd.Nick, gd.Tag)
        }
    }

    summary += "\n🏆 Удачи на турнире eTriathlon 2026!"
    return summary
}

// Константы с правилами дисциплин
const (
    rulesBS = "📋 ПРАВИЛА BRAWL STARS:\n" +
        "Формат: 1v1 (Дружеский бой)\n" +
        "Один из игроков создаёт код команды и приглашает другого.\n" +
        "Второй присоединяется по коду или через приглашение в друзья.\n" +
        "Один из участников создаёт пустую карту в режиме \"Награда за поимку\".\n" +
        "Игроки по очереди выбирают персонажей.\n" +
        "Победителем считается тот, кто выиграл 2 матча.\n" +
        "При счёте 1:1 выбирают персонажа, предложенного судьями."

    rulesCR = "📋 ПРАВИЛА CLASH ROYALE:\n" +
        "Формат: 1v1 (Дружеский бой)\n" +
        "Один из игроков отправляет запрос «Дружеский бой».\n" +
        "Оба игрока должны добавить друг друга в друзья.\n" +
        "Матч проводится до одной победы/ничьи на групповом этапе.\n" +
        "В плей-офф — до одной победы."

    rulesCH = "📋 ПРАВИЛА ШАХМАТ:\n" +
        "Платформа: Chess.com\n" +
        "Контроль времени: 10+3 минуты\n" +
        "Создатель матча выставляет параметры.\n" +
        "Второй игрок получает приглашение или ссылку.\n" +
        "Матч на групповом этапе до одной победы/ничьи.\n" +
        "В плей-офф — до одной победы."
)