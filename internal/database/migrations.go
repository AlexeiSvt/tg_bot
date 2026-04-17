package database

import (
	"database/sql"
	"log"
)

func migrate(db *sql.DB) error {
	_, err := db.Exec(DB_Create)
	if err != nil {
		log.Printf("migrate error: %v", err)
	}
	return err
}
