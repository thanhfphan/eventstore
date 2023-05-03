package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

type Config struct {
	LogLevel    string `env:"LOG_LEVEL"`
	Environment string `env:"ENVIRONMENT"`
	ServiceName string `env:"SERVICE_NAME"`
	HTTPPort    string `env:"HTTP_PORT,default=4012"`
	DB          *DB
}

type DB struct {
	Host          string `env:"DB_HOST"`
	Name          string `env:"DB_NAME"`
	User          string `env:"DB_USER"`
	Password      string `env:"DB_PASSWORD" json:"-"` // zap ignore
	Port          string `env:"DB_PORT"`
	ConnectionURL string `env:"CONNECTION_URL" json:"-"` // zap ignore
}

func LoadFromEnv(ctx context.Context) (*Config, error) {
	log := logging.FromContext(ctx)
	log.Infof("starting load config from ENV ...")

	cfg := &Config{}
	if err := envconfig.ProcessWith(ctx, cfg, envconfig.OsLookuper()); err != nil {
		return nil, errors.New("envconfig.ProcessWith has err=%v", err)
	}

	log.Infow("service cfg", "config", cfg)
	return cfg, nil
}
