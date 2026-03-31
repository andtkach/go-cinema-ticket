package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type movieHandler interface {
	ListMovies(http.ResponseWriter, *http.Request)
	GetMovie(http.ResponseWriter, *http.Request)
	CreateMovie(http.ResponseWriter, *http.Request)
	UpdateMovie(http.ResponseWriter, *http.Request)
	DeleteMovie(http.ResponseWriter, *http.Request)
}

type bookingHandler interface {
	ListBookingAudit(http.ResponseWriter, *http.Request)
	ListSeats(http.ResponseWriter, *http.Request)
	HoldSeat(http.ResponseWriter, *http.Request)
	ConfirmSession(http.ResponseWriter, *http.Request)
	ReleaseSession(http.ResponseWriter, *http.Request)
}

func setupRouter(movieHandler movieHandler, bookingHandler bookingHandler, requireAuth, requireAdmin func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			movieHandler.ListMovies(w, r)
		case http.MethodPost:
			requireAdmin(http.HandlerFunc(movieHandler.CreateMovie)).ServeHTTP(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/movies/{movieID}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			movieHandler.GetMovie(w, r)
		case http.MethodPut:
			requireAdmin(http.HandlerFunc(movieHandler.UpdateMovie)).ServeHTTP(w, r)
		case http.MethodDelete:
			requireAdmin(http.HandlerFunc(movieHandler.DeleteMovie)).ServeHTTP(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	mux.Handle("GET /bookings", requireAdmin(http.HandlerFunc(bookingHandler.ListBookingAudit)))
	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.Handle("POST /movies/{movieID}/seats/{seatID}/hold", requireAuth(http.HandlerFunc(bookingHandler.HoldSeat)))
	mux.Handle("PUT /sessions/{sessionID}/confirm", requireAuth(http.HandlerFunc(bookingHandler.ConfirmSession)))
	mux.Handle("DELETE /sessions/{sessionID}", requireAuth(http.HandlerFunc(bookingHandler.ReleaseSession)))

	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		payload := map[string]string{
			"server":  "go-cinema-ticket",
			"version": getServerVersion(),
			"time":    time.Now().UTC().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload)
	})

	mux.Handle("/", spaHandler("static"))

	return mux
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("→ %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func spaHandler(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := http.Dir(dir).Open(r.URL.Path); err != nil {
			http.ServeFile(w, r, dir+"/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})
}
