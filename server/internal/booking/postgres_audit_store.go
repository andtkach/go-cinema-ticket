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

func (s *postgresAuditStore) ListAll() ([]BookingAudit, error) {
	rows, err := s.db.Query(
		`SELECT b.booked_at, b.movie_id, COALESCE(m.title, ''), b.seat_id, b.user_name
		 FROM bookings b
		 LEFT JOIN movies m ON m.id::text = b.movie_id
		 ORDER BY b.booked_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := []BookingAudit{}
	for rows.Next() {
		var b BookingAudit
		if err := rows.Scan(&b.BookedAt, &b.MovieID, &b.MovieTitle, &b.SeatID, &b.UserName); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}
