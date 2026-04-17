package backup

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"tgbot/internal/models"
	"time"
)

func ExportToCSV(db *sql.DB, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%s: %w", CreateFile, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{HeaderBackup})
	writer.Write([]string{fmt.Sprintf("%s: %s", Generated, time.Now().Format("2006-01-02 15:04:05"))})
	writer.Write([]string{""})

	writer.Write([]string{TableUsers})
	writer.Write([]string{ID, TelegramID, Name, Surname, Class, Disciplines})

	rows, err := db.Query(SelectQuery)
	if err != nil {
		return fmt.Errorf("%s: %w",QueryUsers, err)
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
	writer.Write([]string{fmt.Sprintf("%s: %d", TotalRegistrations, rowCount)})

	writer.Write([]string{"", StatisticsByDiscipline})
	stats, _ := getStatistics(db)
	for disc, count := range stats {
		writer.Write([]string{"", disc, fmt.Sprintf("%d %s", count, Participants)})
	}

	return nil
}
