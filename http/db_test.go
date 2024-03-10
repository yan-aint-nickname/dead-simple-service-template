package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
	Endpoint  string
}

func newPostgresContainer(l *fxtest.Lifecycle) (*PostgresContainer, error) {
	ctx := context.Background()
	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "postgres"

	postgresC, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

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
		Endpoint:  endpoint,
	}, nil
}

func newTestDBPool(l *fxtest.Lifecycle, s SettingsHttp) (*Postgres, error) {
	client, err := NewPostgresPool(l, s)
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

func testDBMigrationUp(t fxtest.TB, pc *PostgresContainer) {
	exe, err := exec.LookPath("goose")
	if err != nil {
		t.Errorf("goose is missing")
	}
	pwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory %s", err)
	}
	projectDir := filepath.Dir(pwd)
	migrationDir := "migrations"

	cmd := exec.Command(
		exe,
		fmt.Sprintf("--dir=%s", filepath.Join(projectDir, migrationDir)),
		"postgres",
		pc.Endpoint,
		"up",
	) // #nosec: G204

	// NOTE: All that thing with stdout and stderr is to catch errors of goose
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		t.Errorf("Postgres migration up test failed %s", err)
	}
	t.Logf(stdBuffer.String())
}

func newPostgresTestOption() fx.Option {
	return fx.Options(
		fx.Provide(
			newPostgresContainer,
			func() *RedisContainer { return nil },
			NewTestSettingsHttp,
			newTestDBPool,
		),
	)
}

func TestPostgres(t *testing.T) {
	t.Run(
		"Ping",
		func(t *testing.T) {
			NewTestApp(t, newPostgresTestOption(), fx.Invoke(testDBPoolConnection)).Stop()
		},
	)
	t.Run(
		"Migrate UP",
		func(t *testing.T) {
			NewTestApp(t, newPostgresTestOption(), fx.Invoke(testDBMigrationUp)).Stop()
		},
	)
}
