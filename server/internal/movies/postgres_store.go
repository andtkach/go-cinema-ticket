package movies

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type postgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) MovieStore {
	return &postgresStore{db}
}

func (s *postgresStore) List() ([]Movie, error) {
	rows, err := s.db.Query(
		`SELECT id, title, rows, seats_per_row, created_at, updated_at
		 FROM movies ORDER BY created_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Rows, &m.SeatsPerRow, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, rows.Err()
}

func (s *postgresStore) GetByID(id string) (Movie, error) {
	var m Movie
	err := s.db.QueryRow(
		`SELECT id, title, rows, seats_per_row, created_at, updated_at FROM movies WHERE id = $1`, id,
	).Scan(&m.ID, &m.Title, &m.Rows, &m.SeatsPerRow, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Movie{}, ErrMovieNotFound
	}
	return m, err
}

func (s *postgresStore) Create(m Movie) (Movie, error) {
	err := s.db.QueryRow(
		`INSERT INTO movies (id, title, rows, seats_per_row)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, title, rows, seats_per_row, created_at, updated_at`,
		m.ID, m.Title, m.Rows, m.SeatsPerRow,
	).Scan(&m.ID, &m.Title, &m.Rows, &m.SeatsPerRow, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return Movie{}, ErrMovieIDConflict
		}
		return Movie{}, err
	}
	return m, nil
}

func (s *postgresStore) Update(m Movie) (Movie, error) {
	err := s.db.QueryRow(
		`UPDATE movies
		 SET title=$2, rows=$3, seats_per_row=$4, updated_at=now()
		 WHERE id=$1
		 RETURNING id, title, rows, seats_per_row, created_at, updated_at`,
		m.ID, m.Title, m.Rows, m.SeatsPerRow,
	).Scan(&m.ID, &m.Title, &m.Rows, &m.SeatsPerRow, &m.CreatedAt, &m.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Movie{}, ErrMovieNotFound
	}
	return m, err
}

func (s *postgresStore) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM movies WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrMovieNotFound
	}
	return nil
}
