package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/eventstore/domain/aggregates"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

type GetOrderRequest struct {
	ID string `json:"id" uri:"id"`
}

func (a *app) handleGetOrder() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Request.Context()
		log := logging.FromContext(ctx)
		log.Infof("starting handle get order ...")

		var req GetOrderRequest
		err := ginCtx.BindUri(&req)
		if err != nil {
			a.handleError(ginCtx, err)
			return
		}

		log.Debugf("handle get order with request info=%v", req)

		var agg aggregates.OrderAggregate
		err = a.aggStore.Get(ctx, req.ID, &agg)
		if err != nil {
			a.handleError(ginCtx, err)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"data": agg,
		})
		return
	}
}
