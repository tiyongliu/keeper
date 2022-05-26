package routers

import (
	"dbbox/app/src/controllers"
	"github.com/gin-gonic/gin"
)

func publicConnections(router *gin.RouterGroup) {
	router.POST("test", controllers.Test)
	router.GET("list")
}
