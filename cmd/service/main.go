package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/thanhfphan/eventstore/app"
	"github.com/thanhfphan/eventstore/config"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/logging"
	"github.com/thanhfphan/eventstore/pkg/server"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	log := logging.NewLogger(os.Getenv("LOG_LEVEL"), os.Getenv("ENVIRONMENT"))
	ctx = logging.WithLogger(ctx, log)

	defer func() {
		done()
		if r := recover(); r != nil {
			log.Errorf("apllication went wrong. Panic err=%v", r)
		}
	}()

	err := realMain(ctx)
	done()
	if err != nil {
		log.Errorf("realMain has failed with err=%v", err)
		return
	}
	log.Infof("APP shutdown successful")
}

func realMain(ctx context.Context) error {
	log := logging.FromContext(ctx)

	cfg, err := config.LoadFromEnv(ctx)
	if err != nil {
		return errors.New("load config from environment failed with err=%v", err)
	}

	// FIXME: move to other places
	db, err := sql.Open("postgres", cfg.DB.ConnectionURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	app, err := app.New(cfg, db)
	if err != nil {
		return errors.New("init server http handler failed with err=%v", err)
	}

	srv, err := server.New(cfg.HTTPPort)
	if err != nil {
		return err
	}

	log.Infof("HTTP Server running on PORT: %s", cfg.HTTPPort)

	return srv.ServeHTTPHandler(ctx, app.Routes(ctx))
}
