package backup

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"tgbot/internal/constants"
	"tgbot/internal/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBackupRoutine(bot *tgbotapi.BotAPI, db *sql.DB) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	msg := tgbotapi.NewMessage(constants.AdminChatID, StartSystemBackup)
	bot.Send(msg)

	for range ticker.C {
		PerformBackup(bot, db)
	}
}

func PerformBackup(bot *tgbotapi.BotAPI, db *sql.DB) {
	filename := fmt.Sprintf("backup_etriathlon_%s.csv", time.Now().Format("2006-01-02_15-04-05"))

	err := ExportToCSV(db, filename)
	if err != nil {
		log.Printf("%s, %v", BackupCreationMistake, err)
		msg := tgbotapi.NewMessage(constants.AdminChatID, fmt.Sprintf("%s: %v", BackupCreationMistake, err))	
		bot.Send(msg)
		return
	}
	defer os.Remove(filename)

	err = SendBackupFile(bot, filename)
	if err != nil {
		log.Printf("%s, %v", BackupSendingMistake, err)
		msg := tgbotapi.NewMessage(constants.AdminChatID, fmt.Sprintf("%s: %v", BackupSendingMistake, err))	
		bot.Send(msg)
		return
	}

	log.Printf("%s, %s", BackupWasSentSuccessfully, filename)
}

func ExportToCSV(db *sql.DB, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"=== eTriathlon 2025 - Database Backup ==="})
	writer.Write([]string{fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05"))})
	writer.Write([]string{""})

	writer.Write([]string{"=== TABLE: users ==="})
	writer.Write([]string{"ID", "Telegram ID", "Имя", "Фамилия", "Класс", "Дисциплины"})

	rows, err := db.Query("SELECT id, tg_id, first_name, last_name, class, disciplines FROM users ORDER BY id")
	if err != nil {
		return fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	rowCount := 0
	for rows.Next() {
		var id, tgID int64
		var firstName, lastName, class string
		var disciplinesJSON []byte

		if err := rows.Scan(&id, &tgID, &firstName, &lastName, &class, &disciplinesJSON); err != nil {
			continue
		}

		var disciplines map[string]models.GameData
		disciplinesStr := string(disciplinesJSON)
		if len(disciplinesJSON) > 0 {
			if err := json.Unmarshal(disciplinesJSON, &disciplines); err == nil {
				disciplinesStr = formatDisciplines(disciplines)
			}
		}

		writer.Write([]string{
			fmt.Sprintf("%d", id),
			fmt.Sprintf("%d", tgID),
			firstName,
			lastName,
			class,
			disciplinesStr,
		})
		rowCount++
	}

	writer.Write([]string{""})
	writer.Write([]string{fmt.Sprintf("Total registrations: %d", rowCount)})

	// Статистика
	writer.Write([]string{"", "=== STATISTICS BY DISCIPLINE ==="})
	stats, _ := getStatistics(db)
	for disc, count := range stats {
		writer.Write([]string{"", disc, fmt.Sprintf("%d участников", count)})
	}

	return nil
}

// Твоя исходная функция (оставлена без изменений, как ты просил)
func SendBackupFile(bot *tgbotapi.BotAPI, filename string) error {
	file := tgbotapi.NewDocument(constants.AdminChatID, tgbotapi.FilePath(filename))

	fileInfo, _ := os.Stat(filename)
	fileSize := float64(fileInfo.Size()) / 1024

	file.Caption = fmt.Sprintf(
		"🔄 Автоматический бэкап базы данных\n"+
			"⏰ %s\n"+
			"📊 Файл: %s\n"+
			"💾 Размер: %.2f KB",
		time.Now().Format("02.01.2006 15:04:05"),
		filename,
		fileSize,
	)

	_, err := bot.Send(file)
	return err
}

// Вспомогательные функции (неэкспортируемые)
func formatDisciplines(disciplines map[string]models.GameData) string {
	result := ""
	for game, data := range disciplines {
		if result != "" {
			result += "; "
		}
		if game == "Chess" {
			result += fmt.Sprintf("%s: %s", game, data.Nick)
		} else {
			result += fmt.Sprintf("%s: %s %s", game, data.Nick, data.Tag)
		}
	}
	return result
}

func getStatistics(db *sql.DB) (map[string]int, error) {
	stats := make(map[string]int)
	rows, err := db.Query("SELECT disciplines FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var disciplinesJSON []byte
		if err := rows.Scan(&disciplinesJSON); err == nil {
			var disciplines map[string]models.GameData
			if err := json.Unmarshal(disciplinesJSON, &disciplines); err == nil {
				for game := range disciplines {
					stats[game]++
				}
			}
		}
	}
	return stats, nil
}