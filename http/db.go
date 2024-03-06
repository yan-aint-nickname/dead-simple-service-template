package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type Postgres struct {
	Pool *pgxpool.Pool
	Ctx  context.Context
}

func NewPostgresPool(lc fx.Lifecycle, settings SettingsHttp) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(settings.GetPostgersDsn())
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dbpool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			go dbpool.Close()
			return nil
		},
	})
	return &Postgres{
		Pool: dbpool,
		Ctx:  ctx,
	}, nil
}
