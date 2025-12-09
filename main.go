package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"nilpotential/whats-kept-in-time/db"
	"nilpotential/whats-kept-in-time/routes/wallpapers"
)

func main() {
	log := slog.Default()
	ctx := context.Background()
	if err := run(ctx, log); err != nil {
		log.Error("Failed to run server", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(ctx context.Context, log *slog.Logger) error {
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("Failed to open database", slog.Any("error", err))
		return err
	}
	defer dbpool.Close()

	db := db.New(dbpool)

	mux := http.NewServeMux()

	mux.Handle("/", http.RedirectHandler("/wallpapers", http.StatusFound))

	wh := wallpapers.NewHandler(log, db)
	mux.Handle("/wallpapers", wh)
	mux.Handle("/wallpapers/{id}", wh)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return http.ListenAndServe(":8080", mux)
}
