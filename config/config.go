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
