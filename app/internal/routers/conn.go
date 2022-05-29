package routers

import (
	"github.com/gin-gonic/gin"
	"keeper/app/pkg/api"
)

//todo 暂时未使用
func NewConnections() *api.ApiDefinition {
	apiDef := api.NewApiDefinition(api.WithName("connection"), api.WithPrefix("connection"))
	apiDef.WithHandler(api.POST, "", func(context *gin.Context) (api.ApiResponse, error) {
		return nil, nil
	})

	return apiDef
}
