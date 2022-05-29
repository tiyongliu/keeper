package routers

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	return addRouter(gin.Default())
}

func addRouter(router *gin.Engine) *gin.Engine {
	connections := router.Group("connections")
	publicConnections(connections)

	return router
}
