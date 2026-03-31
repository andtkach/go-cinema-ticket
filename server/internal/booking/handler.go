package booking

import (
	"log"
	"net/http"
	"time"

	"github.com/andtkach/cinema/internal/auth"
	"github.com/andtkach/cinema/internal/utils"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) *handler {
	return &handler{svc}
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")
	userID := auth.UserIDFromContext(r.Context())
	username := auth.UsernameFromContext(r.Context())
	if username == "" {
		username = userID
	}

	data := Booking{
		UserID:   userID,
		UserName: username,
		SeatID:   seatID,
		MovieID:  movieID,
	}

	session, err := h.svc.Book(data)
	if err != nil {
		log.Println(err)
		return
	}

	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movieID"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	utils.WriteJSON(w, http.StatusCreated, holdResponse{
		SeatID:    seatID,
		MovieID:   session.MovieID,
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")

	bookings := h.svc.ListBookings(movieID)

	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID:    b.SeatID,
			UserID:    b.UserID,
			Booked:    true,
			Confirmed: b.Status == "confirmed",
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}

func (h *handler) ListBookingAudit(w http.ResponseWriter, r *http.Request) {
	rows, err := h.svc.ListAuditBookings()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if rows == nil {
		rows = []BookingAudit{}
	}
	utils.WriteJSON(w, http.StatusOK, rows)
}

type seatInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	userID := auth.UserIDFromContext(r.Context())

	session, err := h.svc.ConfirmSeat(r.Context(), sessionID, userID)
	if err != nil {
		return
	}

	utils.WriteJSON(w, http.StatusOK, sessionResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    userID,
		Status:    session.Status,
	})
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")
	userID := auth.UserIDFromContext(r.Context())

	err := h.svc.ReleaseSeat(r.Context(), sessionID, userID)
	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
