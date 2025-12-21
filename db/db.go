package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(dbpool *pgxpool.Pool) *DB {
	return &DB{
		Pool: dbpool,
	}
}

type DB struct {
	Pool *pgxpool.Pool
}

func (db *DB) GetWallpaperById(ctx context.Context, id string) (data Wallpaper, err error) {
	err = db.Pool.QueryRow(ctx, "SELECT id, title, version, pictureurl FROM wallpapers WHERE id=$1", id).Scan(&data.ID, &data.Title, &data.Version, &data.PictureURL)
	return data, err
}

func (db *DB) ListWallpapersByVersion(ctx context.Context, version string) ([]Wallpaper, error) {
	rows, err := db.Pool.Query(ctx, "SELECT id, title, version, pictureurl FROM wallpapers WHERE version=$1", version)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (data Wallpaper, err error) {
		err = row.Scan(&data.ID, &data.Title, &data.Version, &data.PictureURL)
		return data, err
	})
}

func (db *DB) ListVersions(ctx context.Context) ([]Version, error) {
	rows, err := db.Pool.Query(ctx, "SELECT id, kr FROM versions")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (data Version, err error) {
		err = row.Scan(&data.ID, &data.KR)
		return data, err
	})

}

type Version struct {
	ID string
	CN string
	EN string
	JP string
	KR string
}

type Wallpaper struct {
	ID             string
	Title          string
	Version        string
	PictureURL     string
	TransformedURL string
}
