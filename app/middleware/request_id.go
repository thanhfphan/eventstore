package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ctxKeyRequestID = ctxKey("request_id")

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			uid, err := uuid.NewRandom()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": http.StatusText(http.StatusInternalServerError),
				})
				return
			}
			requestID = uid.String()
		}

		ctx := withRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func RequestIDFromCtx(ctx context.Context) string {
	v := ctx.Value(ctxKeyRequestID)
	if v == nil {
		return ""
	}

	if val, ok := v.(string); ok {
		return val
	}

	return ""
}

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, id)
}
