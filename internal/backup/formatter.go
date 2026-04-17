package backup

import (
	"fmt"
	"tgbot/internal/models"
	"time"
)

func formatDisciplines(disciplines map[string]models.GameData) string {
	result := ""
	for game, data := range disciplines {
		if result != "" {
			result += "; "
		}
		if game == Chess {
			result += fmt.Sprintf("%s: %s", game, data.Nick)
		} else {
			result += fmt.Sprintf("%s: %s %s", game, data.Nick, data.Tag)
		}
	}
	return result
}

func buildCaption(filename string, sizeKB float64) string {
	return fmt.Sprintf(
		CaptionTemplate,
		CaptionTitle,
		time.Now().Format("02.01.2006 15:04:05"),
		filename,
		sizeKB,
	)
}