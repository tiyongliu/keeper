package variable

import "github.com/gin-gonic/gin"

type Application struct {
	Port int
}

type SystemApplication struct {
	*Application
	RegisterHttpRoute func() *gin.Engine
}
