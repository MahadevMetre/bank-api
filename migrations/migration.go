package migrations

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func ApplyMigrations(db *sql.DB, migrationDir string) error {
	if err := goose.Up(db, migrationDir); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}
