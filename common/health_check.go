package common

import (
	"claps-admin/response"
	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	response.Success(ctx, nil, "Claps apis are running successfully!")
}
