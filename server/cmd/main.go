package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/andtkach/cinema/internal/adapters/redis"
	"github.com/andtkach/cinema/internal/auth"
	"github.com/andtkach/cinema/internal/booking"
	"github.com/andtkach/cinema/internal/utils"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)
	mux.Handle("GET /", spaHandler("static"))

	store := booking.NewRedisStore(redis.NewClient("localhost:16379", "redis", "redis"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	issuerURL := os.Getenv("AUTHENTIK_ISSUER_URL")
	clientID := os.Getenv("AUTHENTIK_CLIENT_ID")

	var requireAuth func(http.Handler) http.Handler

	if issuerURL != "" && clientID != "" {
		authMiddleware, err := auth.NewMiddleware(context.Background(), issuerURL, clientID)
		if err != nil {
			log.Fatalf("auth init: %v", err)
		}
		requireAuth = authMiddleware.RequireAuth
		log.Printf("auth enabled: issuer=%s", issuerURL)
	} else {
		requireAuth = func(next http.Handler) http.Handler { return next }
		log.Println("auth disabled: AUTHENTIK_ISSUER_URL or AUTHENTIK_CLIENT_ID not set")
	}

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.Handle("POST /movies/{movieID}/seats/{seatID}/hold", requireAuth(http.HandlerFunc(bookingHandler.HoldSeat)))
	mux.Handle("PUT /sessions/{sessionID}/confirm", requireAuth(http.HandlerFunc(bookingHandler.ConfirmSession)))
	mux.Handle("DELETE /sessions/{sessionID}", requireAuth(http.HandlerFunc(bookingHandler.ReleaseSession)))

	if err := http.ListenAndServe(":17080", mux); err != nil {
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

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}
