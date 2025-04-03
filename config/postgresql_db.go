package config

import (
	"bankapi/constants"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"
)

var DB *sql.DB

func InitDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		constants.PostgresHost, constants.PostgresPort, constants.PostgresUsername, constants.PostgresPassword, constants.PostgresDatabase)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	DB = db

	// Configure connection pool
	db.SetMaxOpenConns(getEnvIntWithDefault("DB_MAX_OPEN_CONNS", 200))
	db.SetMaxIdleConns(getEnvIntWithDefault("DB_MAX_IDLE_CONNS", 50))
	db.SetConnMaxIdleTime(time.Duration(getEnvIntWithDefault("DB_CONN_MAX_IDLE_TIME", 60)) * time.Second)
	db.SetConnMaxLifetime(time.Duration(getEnvIntWithDefault("DB_CONN_MAX_LIFETIME", 240)) * time.Second)

	return db
}

func GetDB() *sql.DB {
	return DB
}

func getEnvIntWithDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
