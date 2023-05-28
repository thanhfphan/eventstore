package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/thanhfphan/eventstore/app/middleware"
	"github.com/thanhfphan/eventstore/config"
	"github.com/thanhfphan/eventstore/domain/aggregates"
	"github.com/thanhfphan/eventstore/domain/repos"
	"github.com/thanhfphan/eventstore/domain/service"
	"github.com/thanhfphan/eventstore/pkg/ev"
)

var (
	_ App = (*app)(nil)
)

type App interface {
	Routes(ctx context.Context) http.Handler
}

type app struct {
	cfg      *config.Config
	aggStore service.AggregateStore
}

func New(cfg *config.Config, dbPool *pgxpool.Pool) (App, error) {
	s := ev.NewSerializer()
	s.RegisterAggregate(&aggregates.OrderAggregate{})

	repos := repos.New(dbPool, s)
	aggStore := service.NewAggregateStore(repos)

	return &app{
		cfg:      cfg,
		aggStore: aggStore,
	}, nil
}

func (a *app) Routes(ctx context.Context) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SetRequestID())
	r.Use(middleware.SetLogger())

	r.GET("/health", a.handleHealth())
	r.POST("/place_order", a.handlePlaceOrder())
	r.POST("/cancel_order", a.handleCancelOrder())

	return r
}
