package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"nilpotential/whats-kept-in-time/routes"
)

func main() {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", routes.CreateHomeHandler())
	http.HandleFunc("/wallpapers", routes.CreateWallpapersHandler(ctx, dbpool))
	http.HandleFunc("/wallpapers/{id}", routes.CreateWallpapersIdHandler(ctx, dbpool))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
