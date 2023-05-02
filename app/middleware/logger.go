package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhfphan/eventstore/pkg/logging"
)

func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		logger := logging.FromContext(ctx)
		if reqID := RequestIDFromCtx(ctx); reqID != "" {
			logger = logger.With("request_id", reqID)
		}

		ctx = logging.WithLogger(ctx, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
