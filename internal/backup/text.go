package backup

const (
	StartSystemBackup = "✅ Система автоматического бэкапа запущена\n⏰ Интервал: каждые 30 минут"
	BackupCreationMistake = "❌ Ошибка создания бэкапа: "
	BackupSendingMistake = "❌ Ошибка отправки бэкапа:"
	BackupWasSentSuccessfully = "Бэкап успешно отправлен"
)

const (
	HeaderBackup = "=== eTriathlon 2025 - Database Backup ==="
	TableUsers = "=== TABLE: users ==="
	ID = "ID"
	TelegramID = "Telegram ID"
	Name = "Name"
	Surname = "Surname"
	Class = "Class"
	Disciplines = "Disciplines"
)

const (
	SelectQuery = `SELECT id, tg_id, first_name, last_name, class, disciplines FROM users ORDER BY id`
	SelectDisciplineFromUser = `SELECT disciplines FROM users`
)

const (
	StatisticsByDiscipline = "=== STATISTICS BY DISCIPLINE ==="
	TotalRegistrations = "Total Registrations: "
)

const CaptionTemplate = "%s\nTimee: %s\nFile: %s\nSize: %.2f KB"

const (
	CaptionTitle = "Authomatic Database Backup"
	CaptionTime  = "Timee: %s"
	CaptionFile  = "File: %s"
	CaptionSize  = "Size: %.2f KB"
)

const (
	CreateFile = "Create File"
	Generated = "Generated"
	QueryUsers = "Query Users"
	Participants = "Participants"
)

const (
	Chess = "Chess"
)

const (
	InvalidAdminID = "Invalid AdminChatID"
)