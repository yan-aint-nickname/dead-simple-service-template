package main

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"log"
	"testing"
)

func TestValidateApp(t *testing.T) {
	if err := fx.ValidateApp(CreateDefaultApp()); err != nil {
		t.Fatal(err)
	}
}

type RedisContainer struct {
	testcontainers.Container
	Endpoint string
}

func SetupRedisContainer() (*RedisContainer, error) {
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
	log.Printf("Redis container endpoint: %s", endpoint)
	if err != nil {
		return nil, err
	}

	return &RedisContainer{Container: redisC, Endpoint: endpoint}, nil
}

func NewTestSettingsHttp(redis_dsn string) *SettingsHttp {
	return &SettingsHttp{
		RedisDsn: redis_dsn,
		RedisExp: 1000 * 1000 * 1000 * 10, // 10 sec
	}
}

func TestRedisConnection(t *testing.T) {
	ctx := context.Background()

	redisC, err := SetupRedisContainer()

	if err != nil {
		t.Fatal(err)
	}
	redisDsn := fmt.Sprintf("redis://%s", redisC.Endpoint)
	testSettings := NewTestSettingsHttp(redisDsn)

	client, err := NewRedisClient(*testSettings)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Client.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := redisC.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop redis: %s", err)
		}
	})
}

// TODO: add fx.Lifecicle test example
