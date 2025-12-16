package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost string
	ServicePort int

	JWT   JWTConfig
	Redis RedisConfig

	DialTimeout time.Duration
	ReadTimeout time.Duration

	ComputingServiceURL   string
	ComputingServiceToken string
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn time.Duration
}

func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)

	if err != nil {
		return nil, err
	}

	if cfg.JWT.SecretKey == "" {
		cfg.JWT.SecretKey = "default-secret-key-change-in-production"
	}
	if cfg.JWT.ExpiresIn == 0 {
		cfg.JWT.ExpiresIn = 24 * time.Hour
	}

	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 3 * time.Second
	}

	cfgRedis, err := ConfRedis(cfg.DialTimeout, cfg.ReadTimeout)
	if err != nil {
		return nil, err
	}

	cfg.Redis = cfgRedis

	log.Info("config parsed")

	return cfg, nil
}
