package main

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
	Endpoint string
}

func newPostgresContainer(l *fxtest.Lifecycle) (*PostgresContainer, error) {
	ctx := context.Background()
	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "postgres"

	postgresC, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		// TODO: add test data
		// postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
		// postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	endpoint, err := postgresC.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}

	l.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return postgresC.Terminate(ctx)
		},
	})

	return &PostgresContainer{
		Container: postgresC,
		Endpoint: endpoint,
	}, nil
}

func newTestDBPool(l *fxtest.Lifecycle, s *SettingsHttp) (*Postgres, error) {
	client, err := NewPostgresPool(l, *s)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func testDBPoolConnection(t fxtest.TB, client *Postgres) {
	ctx := context.Background()
	if err := client.Pool.Ping(ctx); err != nil {
		t.Errorf("Postgres connection test failed %s", err)
	}
}

func newPostgresTestOption() fx.Option {
	return fx.Options(
		fx.Provide(
			newPostgresContainer,
			func() *RedisContainer {return nil},
			NewTestSettingsHttp,
			newTestDBPool,
		),
	)
}

func registerDBTests() fx.Option {
	return fx.Options(
		fx.Invoke(testDBPoolConnection),
	)
}

func TestPostgresConnection(t *testing.T) {
	t.Run(
		"Ping",
		func (t *testing.T) {
			NewTestApp(t, newPostgresTestOption(), registerDBTests()).Stop()
		},
	)
}
// TODO: add migration tests
