package movies

import (
	"errors"
	"time"
)

var (
	ErrMovieNotFound  = errors.New("movie not found")
	ErrMovieIDConflict = errors.New("movie id already exists")
)

type Movie struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Rows        int       `json:"rows"`
	SeatsPerRow int       `json:"seats_per_row"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MovieStore interface {
	List() ([]Movie, error)
	GetByID(id string) (Movie, error)
	Create(m Movie) (Movie, error)
	Update(m Movie) (Movie, error)
	Delete(id string) error
}
