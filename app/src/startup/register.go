package startup

import (
	"dbbox/app/src/routers"
	"github.com/gin-gonic/gin"
)

func RegisterHttpRoute() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	return routers.InitRouter()
}
