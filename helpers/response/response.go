package response

import (
	"user-svc/model"

	"github.com/gin-gonic/gin"
)

func JSON(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, model.ResponseSuccess{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
