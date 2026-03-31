package main

import (
	"database/sql"
	"log"
)

func seedDefaultMovies(db *sql.DB) {
	defaults := []movieSeed{
		{"01960f13-4ec9-7ad0-ae6e-0a8c329f0901", "Inception", 5, 8},
		{"01960f13-4eca-7f6d-9ab3-b0fe1f99c92a", "Dune Part Two", 4, 6},
	}
	for _, m := range defaults {
		_, err := db.Exec(
			`INSERT INTO movies (id, title, rows, seats)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO NOTHING`,
			m.ID, m.Title, m.Rows, m.Seats,
		)
		if err != nil {
			log.Printf("seed movie %s: %v", m.ID, err)
		}
	}
}
