package wallpapers

import (
	"html/template"
	"log/slog"
	"net/http"

	"nilpotential/whats-kept-in-time/db"
)

var templates = template.Must(template.ParseGlob("./templates/*.html"))

type PageData struct {
	Versions   []db.Version
	Wallpapers []db.Wallpaper
}

func NewHandler(log *slog.Logger, db *db.DB) http.Handler {
	return &Handler{
		Log: log,
		DB:  db,
	}
}

type Handler struct {
	Log *slog.Logger
	DB  *db.DB
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var list []db.Wallpaper
	var err error
	if id := r.PathValue("id"); len(id) > 0 {
		data, err := h.DB.GetWallpaperById(r.Context(), id)
		if err != nil {
			h.Log.Error("Failed to get wallpaper", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = []db.Wallpaper{data}
	} else if version := r.FormValue("version"); len(version) > 0 {
		list, err = h.DB.ListWallpapersByVersion(r.Context(), version)
		if err != nil {
			h.Log.Error("Failed to list wallpapers", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	versions, err := h.DB.ListVersions(r.Context())
	if err != nil {
		h.Log.Error("Failed to list versions", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", PageData{versions, list})
	if err != nil {
		h.Log.Error("Failed to execute template", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
