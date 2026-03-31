package movies

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrMovieNotFound  = errors.New("movie not found")
	ErrMovieIDConflict = errors.New("movie id already exists")
)

type Movie struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Rows      int       `json:"rows"`
	Seats     int       `json:"seats"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MovieStore interface {
	List() ([]Movie, error)
	GetByID(id uuid.UUID) (Movie, error)
	Create(m Movie) (Movie, error)
	Update(m Movie) (Movie, error)
	Delete(id uuid.UUID) error
}
