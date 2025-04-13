package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Logger(router *gin.Engine) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		latencyMs := float64(param.Latency.Nanoseconds()) / 1e6
		return fmt.Sprintf("[%s] | %d | %.2fms | %s | %s %s | %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.StatusCode,
			latencyMs,
			param.ClientIP,
			param.Method, param.Path,
			param.Request.UserAgent(),
		)
	})
}
