package main

import (
	"fmt"
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"testing"
)

type RedisContainer struct {
	testcontainers.Container
	Endpoint string
}

func newRedisContainer(l *fxtest.Lifecycle) (*RedisContainer, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:6-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}
	l.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return redisC.Terminate(ctx)
		},
	})

	return &RedisContainer{Container: redisC, Endpoint: endpoint}, nil
}

func newTestSettingsHttp(rc *RedisContainer) (*SettingsHttp, error) {
	endpoint := fmt.Sprintf("redis://%s", rc.Endpoint)
	return &SettingsHttp{RedisDsn: endpoint}, nil
}

func newTestRedisClient(s *SettingsHttp) (*RedisClient, error) {
	client, err := NewRedisClient(*s)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func testRedisConnection(t fxtest.TB, client *RedisClient) {
	ctx := context.Background()

	if err := client.Client.Ping(ctx).Err(); err != nil {
		t.Errorf("Redis connection test failed %s", err)
	}
}

func TestCacheProvide(t *testing.T) {
	app := fxtest.New(t,
		fx.Supply(fx.Annotate(t, fx.As(new(fxtest.TB)))),
		fx.Provide(
			fxtest.NewLifecycle,
			newRedisContainer,
			newTestSettingsHttp,
			newTestRedisClient,
		),
		fx.Invoke(testRedisConnection),
	)

	if err := app.Err(); err != nil {
		t.Fatal(err)
	}
}
