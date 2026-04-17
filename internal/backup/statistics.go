package backup

import (
	"database/sql"
	"encoding/json"
	"tgbot/internal/models"
)

func getStatistics(db *sql.DB) (map[string]int, error) {
	stats := make(map[string]int)
	rows, err := db.Query(SelectDisciplineFromUser)
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
