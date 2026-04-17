package database

import (
	"database/sql"
	"encoding/json"
	"tgbot/internal/models"
)

func SaveUser(db *sql.DB, u *models.User) error {
	disciplinesJSON, err := json.Marshal(u.Disciplines)
	if err != nil {
		return err
	}

	err = db.QueryRow(Insert_Query, u.TelegramID, u.FirstName, u.LastName, u.Class, disciplinesJSON).Scan(&u.ID)

	return err
}