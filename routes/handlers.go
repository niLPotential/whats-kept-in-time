package routes

import (
	"context"
	"html/template"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"nilpotential/whats-kept-in-time/db"
)

var templates = template.Must(template.ParseGlob("./templates/*.html"))

type PageData struct {
	Versions   []db.VersionData
	Wallpapers []db.WallpaperData
}

func CreateHomeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/wallpapers", http.StatusMovedPermanently)
	}
}

func CreateWallpapersHandler(ctx context.Context, dbpool *pgxpool.Pool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		version := r.FormValue("version")
		var list []db.WallpaperData
		var err error
		if len(version) > 0 {
			list, err = db.QueryWallpapersByVersion(ctx, dbpool, version)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		versions, err := db.QueryVersions(ctx, dbpool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = templates.ExecuteTemplate(w, "index.html", PageData{Versions: versions, Wallpapers: list})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CreateWallpapersIdHandler(ctx context.Context, dbpool *pgxpool.Pool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		data, err := db.QueryWallpaperById(ctx, dbpool, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		versions, err := db.QueryVersions(ctx, dbpool)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = templates.ExecuteTemplate(w, "index.html", PageData{Versions: versions, Wallpapers: []db.WallpaperData{data}})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
