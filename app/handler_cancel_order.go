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
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Debugf("handle cancel order with request info=%v", req)

		agg := aggregates.OrderAggregate{}
		err = a.aggStore.Get(ctx, req.OrderID, &agg)
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		log.Debugf("AGGREGATE root: %v", agg.Root())
		log.Debugf("AGGREGATE order: %v", agg.Order)

		agg.RecordCancelled(time.Now())
		err = a.aggStore.Save(ctx, &agg)
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		return
	}
}
