package main

import (
	"context"
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
	// DB setup
	dsn := os.Getenv("POSTGRES_APP_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:15432/cinema?sslmode=disable"
	}
	db := pgadapter.NewClient(dsn)
	seedDefaultMovies(db)

	movieHandler := movies.NewHandler(movies.NewService(movies.NewPostgresStore(db)))

	bookingHandler := booking.NewHandler(
		booking.NewService(
			booking.NewRedisStore(redis.NewClient("localhost:16379", "redis", "redis")),
			booking.NewPostgresAuditStore(db),
		),
	)

	// Auth middleware (optional)
	requireAuth, requireAdmin := defaultAuthMiddleware()
	issuerURL := os.Getenv("AUTHENTIK_ISSUER_URL")
	clientID := os.Getenv("AUTHENTIK_CLIENT_ID")
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

	// Router
	mux := setupRouter(movieHandler, bookingHandler, requireAuth, requireAdmin)
	logged := loggingMiddleware(mux)

	log.Println("server listening on :17080")
	if err := http.ListenAndServe(":17080", logged); err != nil {
		log.Fatal(err)
	}
}

func defaultAuthMiddleware() (func(http.Handler) http.Handler, func(http.Handler) http.Handler) {
	noop := func(next http.Handler) http.Handler { return next }
	return noop, noop
}
