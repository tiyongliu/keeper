package routers

import "github.com/gin-gonic/gin"

func publicConnections(router *gin.RouterGroup) {
	router.GET("test")
	router.GET("list")
}
