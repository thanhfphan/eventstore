package main

import (
	"context"
	stdErr "errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thanhfphan/eventstore/config"
	"github.com/thanhfphan/eventstore/pkg/errors"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

var (
	pathMigration    = flag.String("path", "migrations/", "path to migrations folder")
	migrationTimeout = flag.Duration("timeout", 10*time.Minute, "duration for migration timeout")
)

func main() {
	flag.Parse()

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
	log.Infof("starting migration ...")

	cfg, err := config.LoadFromEnv(ctx)
	if err != nil {
		return errors.New("load config from environment failed with err=%v", err)
	}

	dir := fmt.Sprintf("file://%s", *pathMigration)
	m, err := migrate.New(dir, cfg.DB.ConnectionURL)
	if err != nil {
		return err
	}
	defer m.Close()

	m.LockTimeout = *migrationTimeout
	m.Log = newLogger()

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return errors.New("get migrate version failed with err=%v", err)
	}

	if dirty {
		// FIXME: not work with version 0 :D
		m.Force(int(version) - 1)
	}

	if err := m.Up(); err != nil && !stdErr.Is(err, migrate.ErrNoChange) {
		return errors.New("failed run migrate: %v", err)
	}

	log.Infof("Migration done ...")

	return nil
}

type logger struct {
	logger *log.Logger
}

func newLogger() *logger {
	return &logger{
		logger: log.New(os.Stdout, "migrate", log.LstdFlags),
	}
}

func (l *logger) Printf(arg string, vars ...interface{}) {
	l.logger.Printf(arg, vars...)
}

func (l *logger) Verbose() bool {
	return true
}
