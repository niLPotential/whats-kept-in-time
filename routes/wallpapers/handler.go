package wallpapers

import (
	"html/template"
	"log/slog"
	"net/http"

	"nilpotential/whats-kept-in-time/db"

	"github.com/starfederation/datastar-go/datastar"
)

var templates = template.Must(template.ParseGlob("./templates/*.html"))

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

type Signals struct {
	Version string `json:"version"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if id := r.PathValue("id"); len(id) > 0 {
		wallpaper, err := h.DB.GetWallpaperById(r.Context(), id)
		if err != nil {
			h.Log.Error("Failed to get wallpaper", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		wallpaper.TransformedURL = db.BuildURLFromId(wallpaper.ID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = templates.ExecuteTemplate(w, "modal", wallpaper)
		if err != nil {
			h.Log.Error("Failed to execute template", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	} else if data := r.FormValue("datastar"); len(data) > 0 {
		signals := &Signals{}
		err := datastar.ReadSignals(r, signals)
		if err != nil {
			h.Log.Error("Failed to read datastar signal", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		list, err := h.DB.ListWallpapersByVersion(r.Context(), signals.Version)
		if err != nil {
			h.Log.Error("Failed to list wallpapers", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for i, wallpaper := range list {
			list[i].TransformedURL = db.BuildURLFromId(wallpaper.ID)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = templates.ExecuteTemplate(w, "gallery", list)
		if err != nil {
			h.Log.Error("Failed to execute template", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	versions, err := h.DB.ListVersions(r.Context())
	if err != nil {
		h.Log.Error("Failed to list versions", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", versions)
	if err != nil {
		h.Log.Error("Failed to execute template", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
