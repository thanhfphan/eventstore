package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		log := logging.FromContext(ctx)
		if reqID := RequestIDFromCtx(ctx); reqID != "" {
			log = log.With("request_id", reqID)
		}

		ctx = logging.WithLogger(ctx, log)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
