package integration

import (
	"bankapi/migrations"
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func InitTestDatabase(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
		).WithDeadline(time.Minute * 2),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Cleanup function to terminate the container
	cleanup := func() {
		err := postgresC.Terminate(ctx)
		if err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}

	// Get the container's host and port
	host, err := postgresC.Host(ctx)
	if err != nil {
		cleanup()
		t.Fatal(err)
	}
	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		cleanup()
		t.Fatal(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=postgres password=testpassword dbname=testdb sslmode=disable", host, port.Int())

	var db *sql.DB

	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect, retrying in 2 seconds... (attempt %d/3)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err := migrations.ApplyMigrations(db, "../../migrations"); err != nil {
		t.Fatalf("migration error: %v", err)
	}

	if err != nil {
		cleanup()
		t.Fatalf("Failed to connect to the database after 3 attempts: %v", err)
	}

	return db, cleanup
}
