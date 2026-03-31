package booking

import "context"

type Service struct {
	store      BookingStore
	auditStore BookingAuditStore
}

func NewService(store BookingStore, auditStore ...BookingAuditStore) *Service {
	var as BookingAuditStore
	if len(auditStore) > 0 {
		as = auditStore[0]
	}
	return &Service{store: store, auditStore: as}
}

func (s *Service) Book(b Booking) (Booking, error) {
	session, err := s.store.Book(b)
	if err != nil {
		return Booking{}, err
	}
	if s.auditStore != nil {
		if err := s.auditStore.InsertHold(Booking{
			MovieID:  session.MovieID,
			SeatID:   session.SeatID,
			UserName: b.UserName,
		}); err != nil {
			return Booking{}, err
		}
	}
	return session, nil
}

func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}

func (s *Service) ListAuditBookings() ([]BookingAudit, error) {
	if s.auditStore == nil {
		return []BookingAudit{}, nil
	}
	return s.auditStore.ListAll()
}

func (s *Service) ConfirmSeat(ctx context.Context, sessionID string, userID string) (Booking, error) {
	return s.store.Confirm(ctx, sessionID, userID)
}

func (s *Service) ReleaseSeat(ctx context.Context, sessionID string, userID string) error {
	return s.store.Release(ctx, sessionID, userID)
}
