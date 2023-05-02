package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/eventstore/app/middleware"
	"github.com/thanhfphan/eventstore/config"
)

var (
	_ App = (*app)(nil)
)

type App interface {
	Routes(ctx context.Context) http.Handler
}

type app struct {
	cfg *config.Config
}

func New(cfg *config.Config) (App, error) {
	return &app{
		cfg: cfg,
	}, nil
}

func (a *app) Routes(ctx context.Context) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SetRequestID())
	r.Use(middleware.SetLogger())

	r.GET("/health", a.handleHealth())

	return r
}
