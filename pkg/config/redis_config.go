package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

func ConfRedis(dialTimeout, readTimeout time.Duration) (RedisConfig, error) {
	host := os.Getenv(envRedisHost)
	if host == "" {
		return RedisConfig{}, fmt.Errorf("REDIS_HOST is required")
	}

	portStr := os.Getenv(envRedisPort)
	if portStr == "" {
		return RedisConfig{}, fmt.Errorf("REDIS_PORT is required")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return RedisConfig{}, fmt.Errorf("redis port must be int value: %w", err)
	}

	user := os.Getenv(envRedisUser)
	pass := os.Getenv(envRedisPass)

	return RedisConfig{
		Host:        host,
		Port:        port,
		User:        user,
		Password:    pass,
		DialTimeout: dialTimeout,
		ReadTimeout: readTimeout,
	}, nil
}
