package routers

import (
	"github.com/gin-gonic/gin"
	"keeper/app/internal/controllers"
)

func publicConnections(router *gin.RouterGroup) {
	router.POST("test", controllers.Test)
	router.GET("list")
}
