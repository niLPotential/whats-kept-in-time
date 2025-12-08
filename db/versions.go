package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VersionData struct {
	Id string
	Cn string
	En string
	Jp string
	Kr string
}

func QueryVersions(ctx context.Context, dbpool *pgxpool.Pool) ([]VersionData, error) {
	rows, err := dbpool.Query(ctx, "SELECT id, kr from versions")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (VersionData, error) {
		var data VersionData
		err := row.Scan(&data.Id, &data.Kr)
		return data, err
	})
}
