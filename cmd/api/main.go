package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/teddy-codes/url-shortner/internal/services/url"
	"github.com/teddy-codes/url-shortner/internal/store/postgres"
)

func middlewareLogger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "logger", log)
			r = r.WithContext(ctx)

			log.Info().
				Str("remote_ip", r.RemoteAddr).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Request Received")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func main() {
	loggerOut := zerolog.NewConsoleWriter()
	log := zerolog.New(loggerOut)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=password dbname=url_shortner sslmode=disable")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot ping database")
	}
	log.Info().Msg("successfully connected to the db")

	store := &postgres.Store{
		DB: db,
	}

	if err := store.CreateURLTable(ctx); err != nil {
		log.Fatal().Err(err).Msg("cannot create url table")
	}
	log.Info().Msg("successfully created URL table")

	svc := &url.Service{
		Store: store,
	}

	h := &handlers{
		log:     log,
		service: svc,
	}

	router := chi.NewMux()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middlewareLogger(log))
	router.Post("/", h.createURL)
	router.Get("/{path}", h.getURL)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Err(err).Msg("cannot serve on port 8080")
	}
}
