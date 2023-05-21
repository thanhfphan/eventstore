package app

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/eventstore/domain/aggregates"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

type PlaceOrderRequest struct {
	CustomerID int64   `json:"customer_id"`
	Price      float64 `json:"price"`
}

func (a *app) handlePlaceOrder() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := ginCtx.Request.Context()
		log := logging.FromContext(ctx)
		log.Infof("starting handle place orderapi ...")

		var req PlaceOrderRequest
		err := ginCtx.Bind(&req)
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		log.Debugf("handle place order with request info=%v", req)

		now := time.Now()
		agg := aggregates.CreateOrderAggregate(req.CustomerID, req.Price, now)
		err = a.aggStore.Save(ctx, agg)
		if err != nil {
			ginCtx.JSON(http.StatusBadRequest, gin.H{
				"err": err,
			})
			return
		}

		ginCtx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
		return
	}
}
