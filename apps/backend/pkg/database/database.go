package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func SetupDb() (*sql.DB, error) {
	config := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432, // Default to 5432 if DB_PORT is not set
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SLL_MODE"),
	}

	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		} else {
			fmt.Printf("Warning: Invalid DB_PORT value '%s', using default 5432\n", portStr)
		}
	}

	fmt.Printf("Attempting to connect to database with config: %+v\n", config)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	err = SetupTables(sqlDB)
	if err != nil {
		return nil, fmt.Errorf("failed to setup tables: %w", err)
	}

	fmt.Println("Successfully connected to the database")

	return sqlDB, nil

}

func SetupTables(db *sql.DB) error {
	stmt := `
CREATE TABLE IF NOT EXISTS "users" (
"id" VARCHAR(128) PRIMARY KEY,
"first_name" VARCHAR(255) NOT NULL,
"last_name" VARCHAR(255) NOT NULL,
"email" VARCHAR(255) UNIQUE NOT NULL,
"password" VARCHAR(255) NOT NULL,
"telephone_number" VARCHAR(11) UNIQUE NOT NULL,
"point" DOUBLE PRECISION,
"created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
"updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

`

	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	stmt = `
CREATE TABLE  IF NOT EXISTS"companys" (
      "id" VARCHAR(128) PRIMARY KEY,
      "company_name" VARCHAR(255) NOT NULL,
      "email" VARCHAR(255) UNIQUE NOT NULL,
      "password" VARCHAR(255) NOT NULL,
      "telephone_number" VARCHAR(11) UNIQUE NOT NULL,
      "logo" VARCHAR(255),
      "point" DOUBLE PRECISION,
      "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
      "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

`

	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create chats table: %w", err)
	}

	stmt = `
CREATE TABLE IF NOT EXISTS "work_detail" (
"id" UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
"start_date" TIMESTAMP NOT NULL,
"end_date" TIMESTAMP NOT NULL,
"trailer_type" VARCHAR(255) NOT NULL,
"weight" VARCHAR(100) NOT NULL,
"description" TEXT,
"completed" VARCHAR(1) NOT NULL
);

`

	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create chats table: %w", err)
	}

	stmt = `
CREATE TABLE IF NOT EXISTS "works" (
      "company_id" VARCHAR(128),
      "work_detail_id" UUID,
      "user_id" VARCHAR(128),
      "start_date" VARCHAR(255) NOT NULL,
      "end_date" VARCHAR(255) NOT NULL,
      "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
      "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
      PRIMARY KEY ("work_detail_id", "company_id")
);

ALTER TABLE "works" ADD FOREIGN KEY ("company_id") REFERENCES "companys" ("id");

ALTER TABLE "works" ADD FOREIGN KEY ("work_detail_id") REFERENCES "work_detail" ("id");

ALTER TABLE "works" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


`

	_, err = db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("failed to create chats table: %w", err)
	}

	fmt.Println("tables successfully created")
	return nil
}
