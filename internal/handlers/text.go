package handlers

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

const (
	UseStartMessage = "Use /start to register, /cancel to abort, /mystats to view your data."

	RulesConfirmPrompt = "Нажмите кнопку ниже, если ознакомились с правилами:"

	TriathlonRules = "🏆 ПРАВИЛА ТРИАТЛОНА\n\n" +
		"Вы участвуете во всех трёх играх:\n" +
		"• Brawl Stars\n" +
		"• Clash Royale\n" +
		"• Chess (Шахматы)\n\n" +
		"Для каждой игры необходимо ввести ник и тег игрока.\n\n" +
		"Выберите игру для ввода данных:"

	EnterNick = "Введите ваш ник в %s:"

	TriathlonIncomplete = "❌ Необходимо заполнить данные для всех трёх игр!"

	ChooseNextGame = "Выберите следующую игру:"

	ConfirmationHeader = "📋 ПРОВЕРКА ДАННЫХ\n\n" +
		"Пожалуйста, внимательно проверьте введённую информацию:\n\n"

	ConfirmationBody = "👤 Имя: %s\n" +
		"👤 Фамилия: %s\n" +
		"📚 Класс: %s\n\n" +
		"🎮 Дисциплины:\n"

	ConfirmationFooter = "\n✅ Я подтверждаю правильность введённой информации и согласен с её обработкой в соответствии с правилами турнира eTriathlon 2026.\n\n" +
		"⚠️ Если вы обнаружили ошибку, удалите чат с ботом и создайте новый для корректировки данных."

	ConfirmButton = "✅ Всё верно, завершить"
	CancelButton  = "❌ Отменить"

	SaveError = "❌ Ошибка при сохранении данных. Попробуйте позже."

	RegistrationCancelled = "❌ Регистрация отменена.\n\nДля начала новой регистрации введите /start"

	SummaryHeader = "✅ Регистрация завершена!\n\n" +
		"Ваши данные:\n" +
		"📝 Имя: %s\n" +
		"📝 Фамилия: %s\n" +
		"📚 Класс: %s\n\n" +
		"🎮 Дисциплины:\n"

	SummaryFooter = "\n🏆 Удачи на турнире eTriathlon 2026!"
)

const (

	CBDiscBS  = "disc_bs"
	CBDiscCR  = "disc_cr"
	CBDiscCH  = "disc_ch"
	CBDiscTri = "disc_tri"

	CBTriBS = "tri_bs"
	CBTriCR = "tri_cr"
	CBTriCH = "tri_ch"

	CBTriCheck = "tri_check"
	CBTriDone  = "tri_done"

	CBMoreYes = "more_yes"
	CBMoreNo  = "more_no"

	CBTriConfirm   = "tri_confirm"
	CBFinalConfirm = "final_confirm"

	CBCancel = "cancel_reg"

	CBPrefixOK = "ok_"
)

const (
	MsgAskLastName   = "Введите вашу фамилию:"
	MsgAskClass      = "Введите ваш класс (например: 9А, 10Б):"
	MsgChooseDisc    = "Выберите дисциплину для участия:"
	MsgChessSaved    = "✅ Данные для Chess сохранены!\n\nВыберите следующую игру:"
	MsgMoreGames     = "Хотите зарегистрироваться в других играх?"
	MsgInvalidTag    = "❌ Неверный формат тега! Используйте формат #ABC123"
	MsgNoGame        = "❌ Ошибка: игра не выбрана. Начните заново с /start"
	MsgEnterTag      = "Введите ваш тег в %s (например: #ABC123):"
	MsgChessDone     = "Chess"
)

const (
	TriStatusHeader = "📊 Статус заполнения триатлона:\n\n"
	TriStatusLine   = "%s %s: %s\n"

	TriNotFilled = "не заполнено"
	TriNickOnly  = "ник: %s"
	TriNickTag   = "ник: %s, тег: %s"
)