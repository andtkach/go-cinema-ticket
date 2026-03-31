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
			id            UUID        PRIMARY KEY,
			title         TEXT        NOT NULL,
			rows          INTEGER     NOT NULL CHECK (rows > 0),
			seats         INTEGER     NOT NULL CHECK (seats > 0),
			created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS bookings (
			id         BIGSERIAL   PRIMARY KEY,
			booked_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
			movie_id   TEXT        NOT NULL,
			seat_id    TEXT        NOT NULL,
			user_name  TEXT        NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	return nil
}
