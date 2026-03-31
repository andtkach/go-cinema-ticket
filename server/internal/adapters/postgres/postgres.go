package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewClient(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("postgres: open: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("postgres: ping: %v", err)
	}
	if err := runMigrations(db); err != nil {
		log.Fatalf("postgres: migrations: %v", err)
	}
	log.Println("postgres: connected")
	return db
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
			id            TEXT        PRIMARY KEY,
			title         TEXT        NOT NULL,
			rows          INTEGER     NOT NULL CHECK (rows > 0),
			seats_per_row INTEGER     NOT NULL CHECK (seats_per_row > 0),
			created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	return err
}
