package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	pgadapter "github.com/andtkach/cinema/internal/adapters/postgres"
	"github.com/andtkach/cinema/internal/adapters/redis"
	"github.com/andtkach/cinema/internal/auth"
	"github.com/andtkach/cinema/internal/booking"
	"github.com/andtkach/cinema/internal/movies"
)

func main() {
	mux := http.NewServeMux()

	// Postgres — movies
	dsn := os.Getenv("POSTGRES_APP_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:15432/cinema?sslmode=disable"
	}
	db := pgadapter.NewClient(dsn)
	seedDefaultMovies(db)

	movieStore := movies.NewPostgresStore(db)
	movieSvc := movies.NewService(movieStore)
	movieHandler := movies.NewHandler(movieSvc)

	// Redis — bookings
	store := booking.NewRedisStore(redis.NewClient("localhost:16379", "redis", "redis"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	// Auth middleware
	issuerURL := os.Getenv("AUTHENTIK_ISSUER_URL")
	clientID := os.Getenv("AUTHENTIK_CLIENT_ID")

	requireAuth := func(next http.Handler) http.Handler { return next }
	requireAdmin := func(next http.Handler) http.Handler { return next }

	if issuerURL != "" && clientID != "" {
		authMiddleware, err := auth.NewMiddleware(context.Background(), issuerURL, clientID)
		if err != nil {
			log.Fatalf("auth init: %v", err)
		}
		requireAuth = authMiddleware.RequireAuth
		requireAdmin = func(next http.Handler) http.Handler {
			return authMiddleware.RequireAuth(authMiddleware.RequireGroup("cinema-admins", next))
		}
		log.Printf("auth enabled: issuer=%s", issuerURL)
	} else {
		log.Println("auth disabled: AUTHENTIK_ISSUER_URL or AUTHENTIK_CLIENT_ID not set")
	}

	// Routes — movies (from DB)
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

	// Routes — bookings
	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.Handle("POST /movies/{movieID}/seats/{seatID}/hold", requireAuth(http.HandlerFunc(bookingHandler.HoldSeat)))
	mux.Handle("PUT /sessions/{sessionID}/confirm", requireAuth(http.HandlerFunc(bookingHandler.ConfirmSession)))
	mux.Handle("DELETE /sessions/{sessionID}", requireAuth(http.HandlerFunc(bookingHandler.ReleaseSession)))

	mux.Handle("/", spaHandler("static"))

	logged := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("→ %s %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(":17080", logged); err != nil {
		log.Fatal(err)
	}
}

// spaHandler serves static files and falls back to index.html for unknown
// paths so that the React SPA can handle client-side routing (e.g. /callback).
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

type cinemaResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Screens  int    `json:"screens"`
}

func seedDefaultMovies(db *sql.DB) {
	defaults := []struct {
		id    string
		title string
		rows  int
		seats int
	}{
		{"01960f13-4ec9-7ad0-ae6e-0a8c329f0901", "Inception", 5, 8},
		{"01960f13-4eca-7f6d-9ab3-b0fe1f99c92a", "Dune: Part Two", 4, 6},
	}
	for _, m := range defaults {
		_, err := db.Exec(
			`INSERT INTO movies (id, title, rows, seats)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO NOTHING`,
			m.id, m.title, m.rows, m.seats,
		)
		if err != nil {
			log.Printf("seed movie %s: %v", m.id, err)
		}
	}
}
