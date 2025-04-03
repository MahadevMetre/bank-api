package integration

import (
	"fmt"
	"testing"
)

func TestApplyMigrations(t *testing.T) {
	db, cleanup := InitTestDatabase(t)
	defer cleanup()
	defer db.Close()

	t.Run("verify tables", func(t *testing.T) {
		rows, err := db.Query(`SELECT table_name FROM information_schema.tables WHERE table_schema='public'`)
		if err != nil {
			t.Fatalf("Failed to query tables: %v", err)
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				t.Fatalf("Failed to scan table name: %v", err)
			}
			tables = append(tables, tableName)
		}

		if err := rows.Err(); err != nil {
			t.Fatalf("Error occurred during table rows iteration: %v", err)
		}

		fmt.Println("Tables after migration:")
		for _, table := range tables {
			fmt.Println(table)
		}

		if len(tables) == 0 {
			t.Fatalf("No tables found in the database, migration might have failed.")
		}
	})
}
