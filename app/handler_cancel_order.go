package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/eventstore/domain/aggregates"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

type CancelOrderRequest struct {
	OrderID string `json:"order_id"`
}

func (a *app) handleCancelOrder() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Request.Context()
		log := logging.FromContext(ctx)
		log.Infof("starting handle place orderapi ...")

		var req CancelOrderRequest
		err := ginCtx.Bind(&req)
		if err != nil {
			a.handleError(ginCtx, err)
			return
		}

		log.Debugf("handle cancel order with request info=%v", req)

		agg := aggregates.OrderAggregate{}
		err = a.aggStore.Get(ctx, req.OrderID, &agg)
		if err != nil {
			a.handleError(ginCtx, err)
			return
		}

		log.Debugf("AGGREGATE : %v", agg)

		agg.RecordCancelled(time.Now())
		err = a.aggStore.Save(ctx, &agg)
		if err != nil {
			a.handleError(ginCtx, err)
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		return
	}
}
