package booking

import "database/sql"

type postgresAuditStore struct {
	db *sql.DB
}

func NewPostgresAuditStore(db *sql.DB) BookingAuditStore {
	return &postgresAuditStore{db: db}
}

func (s *postgresAuditStore) InsertHold(b Booking) error {
	_, err := s.db.Exec(
		`INSERT INTO bookings (movie_id, seat_id, user_name) VALUES ($1, $2, $3)`,
		b.MovieID, b.SeatID, b.UserName,
	)
	return err
}
