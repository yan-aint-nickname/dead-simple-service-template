package main

import (
	"fmt"
	"time"
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx context.Context
	Exp time.Duration
}

type IRedis interface {
	Get() ([]byte, error)
	Set() error
}

func NewRedisClient(s SettingsHttp) (*RedisClient, error) {
	redisDsn := s.GetRedisDsn()
	opts, err := redis.ParseURL(redisDsn)
	if err != nil {
		return &RedisClient{}, err
	}
	// NOTE: convert nanoseconds to seconds
	exp := time.Duration(s.RedisExp*1000*1000*1000)
	return &RedisClient{
		Client: redis.NewClient(opts),
		Ctx: context.Background(),
		Exp: exp,
	}, nil
}

func (c RedisClient) Set(key string, value any) error {
	err := c.Client.Set(c.Ctx, key, value, c.Exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c RedisClient) Get(key string) ([]byte, error) {
	res, err := c.Client.Get(c.Ctx, key).Bytes()
	if err == redis.Nil {
		return []byte{}, fmt.Errorf("Key: %s does't exists", key)
	} else if err != nil {
		return []byte{}, err
	}
	return res, nil
}
