package v1

import (
	"github.com/antoniokichaev/hezzl-collector/internal/controller/http/v1/campaign"
	"github.com/antoniokichaev/hezzl-collector/internal/controller/http/v1/item"
	"github.com/antoniokichaev/hezzl-collector/internal/service"
	"github.com/gin-gonic/gin"
)

func New(route *gin.Engine, services *service.Services) {
	baseRoute := route.Group("/v1")

	itemHandlers := item.New(services.Item)
	itemRoute := baseRoute.Group("/item")
	{
		itemRoute.POST("/create", itemHandlers.Create)
		itemRoute.PATCH("/update", itemHandlers.Update)
		itemRoute.DELETE("/remove", itemHandlers.Delete)
		itemRoute.GET("/", itemHandlers.GetItem)
	}
	baseRoute.GET("/items/list", itemHandlers.GetItems)

	campaignHandlers := campaign.New(services.Campaign)
	campaignRoute := baseRoute.Group("/campaign")
	{
		campaignRoute.POST("/create", campaignHandlers.Create)
		campaignRoute.PATCH("/update", campaignHandlers.Update)
		campaignRoute.DELETE("/remove", campaignHandlers.Delete)
	}
	baseRoute.GET("/campaigns/list", campaignHandlers.GetCampaigns)

}
