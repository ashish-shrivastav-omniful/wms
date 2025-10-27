package routes

import (
	"wms/controller"
	"wms/utils"

	"github.com/omniful/go_commons/http"
)

func GetRoutes(router *http.Server) {
	hub := router.Group("/api/hub")
	{
		newHubController := controller.CreateHubController(utils.PostgresDB, &utils.Ctx)
		hub.GET("/view", newHubController.ViewHubs)
		hub.POST("/create", newHubController.CreateHub)
	}
	inventory := router.Group("/api/inventory")
	{
		newInventoryController := controller.CreateInventoryController(utils.PostgresDB,  &utils.Ctx)
		inventory.GET("/view", newInventoryController.ViewInventory)
		inventory.PUT("/edit", newInventoryController.EditInventory)
		inventory.POST("/update", newInventoryController.UpdateInventoryAftersales)
	}
	sku := router.Group("/api/sku")
	{
		newSkuController := controller.CreatSkuController(utils.PostgresDB, &utils.Ctx)
		sku.GET("/view", newSkuController.ViewSkus)
		sku.POST("/create", newSkuController.CreateSkus)
		sku.POST("/verify", newSkuController.VerifySkus)
	}
	seller := router.Group("/api/seller")
	{
		newSellerController := controller.CreateSellerController(utils.PostgresDB,  &utils.Ctx)
		seller.GET("/view", newSellerController.ViewSellers)
		seller.POST("/create", newSellerController.CreateSeller)
	}
	tenant := router.Group("/api/tenant")
	{
		newTenantController := controller.CreateTenantController(utils.PostgresDB,  &utils.Ctx)
		tenant.GET("/view", newTenantController.ViewTenants)
		tenant.POST("/create", newTenantController.CreateTenant)
	}
}