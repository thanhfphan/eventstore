package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thanhfphan/eventstore/config"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/logging"
	"github.com/thanhfphan/eventstore/pkg/server"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger := logging.NewLogger(os.Getenv("LOG_LEVEL"), os.Getenv("ENVIRONMENT"))
	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Errorf("apllication went wrong. Panic err=%v", r)
		}
	}()

	err := realMain(ctx)
	done()
	if err != nil {
		logger.Errorf("realMain has failed with err=%v", err)
		return
	}
	logger.Infof("APP shutdown successful")
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)

	cfg, err := config.LoadFromEnv(ctx)
	if err != nil {
		return errors.New("load config from environment failed with err=%v", err)
	}
	srv, err := server.New(cfg.HTTPPort)
	if err != nil {
		return err
	}

	log.Infof("HTTP Server running on PORT: %s", cfg.HTTPPort)

	return srv.ServeHTTPHandler(ctx, nil)
}
