package startup

import (
	"github.com/gin-gonic/gin"
	"keeper/app/src/routers"
)

func RegisterHttpRoute() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	return routers.InitRouter()
}
