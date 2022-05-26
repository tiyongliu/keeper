package routers

import (
	"dbbox/app/src/nest"
	"github.com/gin-gonic/gin"
)

/*
func NewTuniuApiDefinition() *api.ApiDefinition {
	apiDef := api.NewApiDefinition(api.WithName("Tuniu"), api.WithPrefix("tuniu"))
	tuniuApiServer := tuniu.NewTuniuApiServer(context.Background())

	return apiDef.
		WithHandler(api.GET, "ping", api.Ping).
		WithHandler(api.POST, "search", tuniuApiServer.Search).
		WithHandler(api.POST, "verify", tuniuApiServer.Verify).
		WithEncryHandle(api.POST, "order", config.GetTuniuAesKey(), tuniuApiServer.Order).
		WithHandler(api.POST, "ticketing", tuniuApiServer.Ticketing).
		WithHandler(api.POST, "cancel", tuniuApiServer.Cancel).
		WithHandler(api.POST, "ticketResult", tuniuApiServer.TicketResult)

}

*/

func NewConnections() *nest.ApiDefinition {
	apiDef := nest.NewApiDefinition(nest.WithName("conn"), nest.WithPrefix("conn"))
	apiDef.WithHandler(nest.POST, "", func(context *gin.Context) (nest.ApiResponse, error) {
		return nil, nil
	})

	return apiDef
}
