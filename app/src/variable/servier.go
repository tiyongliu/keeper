package variable

import "github.com/gin-gonic/gin"

type Application struct {
	Port int
}

type SystemApplication struct {
	*Application
	RegisterHandlerRoute func() *gin.Engine
}

