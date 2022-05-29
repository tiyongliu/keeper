package routers

import (
	"github.com/gin-gonic/gin"
	"keeper/app/controllers"
)

func publicConnections(router *gin.RouterGroup) {
	router.POST("test", controllers.Test)
	router.GET("list")
}
