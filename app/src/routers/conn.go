package routers

import (
	"github.com/gin-gonic/gin"
	"keeper/app/src/api"
)

func NewConnections() *api.ApiDefinition {
	apiDef := api.NewApiDefinition(api.WithName("conn"), api.WithPrefix("conn"))
	apiDef.WithHandler(api.POST, "", func(context *gin.Context) (api.ApiResponse, error) {
		return nil, nil
	})

	return apiDef
}
