package pkg

import (
	"context"
	"develop-internet-applications/pkg/config"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	servicePrefix = "exocalc:"
	jwtPrefix     = "jwt:"
)

type RedisClient struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*RedisClient, error) {
	client := &RedisClient{}

	client.cfg = cfg

	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.Password,
		Username:    cfg.User,
		Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	client.client = redisClient

	return client, nil
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func getJWTKey(token string) string {
	return servicePrefix + jwtPrefix + token
}

func (c *RedisClient) WriteJWTToBlacklist(ctx context.Context, jwtStr string, jwtTTL time.Duration) error {
	return c.client.Set(ctx, getJWTKey(jwtStr), true, jwtTTL).Err()
}

func (c *RedisClient) CheckJWTInBlacklist(ctx context.Context, jwtStr string) (bool, error) {
	err := c.client.Get(ctx, getJWTKey(jwtStr)).Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
