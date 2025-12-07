package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var templates = template.Must(template.ParseGlob("./templates/*.html"))

type Data struct {
	Id    string
	Title string
}

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var title string
		err = dbpool.QueryRow(context.Background(), "SELECT title FROM wallpapers WHERE id=$1", id).Scan(&title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = templates.ExecuteTemplate(w, "index.html", Data{Id: id, Title: title})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
