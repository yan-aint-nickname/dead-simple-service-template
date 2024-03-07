package main

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type RedisContainer struct {
	Container *redis.RedisContainer
	Endpoint  string
}

func newRedisContainer(l *fxtest.Lifecycle) (*RedisContainer, error) {
	ctx := context.Background()
	redisC, err := redis.RunContainer(
		ctx,
		testcontainers.WithImage("redis:6-alpine"),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)

	if err != nil {
		return nil, err
	}

	endpoint, err := redisC.ConnectionString(ctx)
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

func newRedisTestOption() fx.Option {
	return fx.Options(
		fx.Provide(
			newRedisContainer,
			func() *PostgresContainer { return nil },
			NewTestSettingsHttp,
			newTestRedisClient,
		),
	)
}

func registerCacheTests() fx.Option {
	return fx.Options(
		fx.Invoke(testRedisConnection),
	)
}

func TestRedisConnection(t *testing.T) {
	t.Run(
		"Ping",
		func(t *testing.T) {
			NewTestApp(t, newRedisTestOption(), registerCacheTests()).Stop()
		},
	)
}
