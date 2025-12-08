package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WallpaperData struct {
	Id         string
	Title      string
	Version    string
	PictureUrl string
}

func QueryWallpaperById(ctx context.Context, dbpool *pgxpool.Pool, id string) (WallpaperData, error) {
	var data WallpaperData
	err := dbpool.QueryRow(ctx, "SELECT id, title, version, pictureurl FROM wallpapers WHERE id=$1", id).Scan(&data.Id, &data.Title, &data.Version, &data.PictureUrl)
	return data, err
}

func QueryWallpapersByVersion(ctx context.Context, dbpool *pgxpool.Pool, version string) ([]WallpaperData, error) {
	rows, err := dbpool.Query(ctx, "SELECT id, title, version, pictureurl from wallpapers WHERE version=$1", version)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (WallpaperData, error) {
		var data WallpaperData
		err := row.Scan(&data.Id, &data.Title, &data.Version, &data.PictureUrl)
		return data, err
	})
}
