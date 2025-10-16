package controller

import (
	"context"
	"net/http"
	"wms/models"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type SkuHandler struct {
	DB  *gorm.DB
	RC  redis.Client
	Ctx *context.Context
}

func CreatSkuHandler(db *gorm.DB, rd redis.Client, ctx *context.Context) *SkuHandler {
	return &SkuHandler{DB: db, RC: rd, Ctx: ctx}
}

// view all skus present
func (h *SkuHandler) viewSkus(c *gin.Context) {
	var skus []models.Sku
	if results := h.DB.Find(&skus); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during fetching skus data",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, skus)
}

// add new skus
func (h *SkuHandler) createSkus(c *gin.Context) {
	var skujson struct {
		TenantId uint   `json:"prodid" binding:"required"`
		Name     string `json:"name" binding:"required"`
		SkuCode  string `json:"skucode" binding:"required"`
		Desc     string `json:"desc" binding:"required"`
		Price    uint   `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&skujson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error binding body data",
			"error":   err.Error(),
		})
		return
	}
	var newSku models.Sku
	newSku.TenantId = skujson.TenantId
	newSku.Name = skujson.Name
	newSku.SkuCode = skujson.SkuCode
	newSku.Desc = skujson.Desc
	newSku.Price = skujson.Price
	if err := h.DB.Save(&newSku); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error saving data in database",
			"error":   err.Error.Error(),
		})
		return 
	}
	c.JSON(200, gin.H{
		"message": "created new sku",
	})
}

// verify from orders to check whether all skus present in a hub or not
func (h *SkuHandler) verifySkus(c *gin.Context) {
	type Order struct {
		SNo      string `json:"sno"`
		SellerID string `json:"seller_id"`
		OrderID  string `json:"order_id"`
		ItemID   string `json:"item_id"`
		Quantity string `json:"quantity"`
		Status   string `json:"status"`
	}

	type OrderResponse struct {
		ValidOrders   []Order `json:"valid_orders"`
		MissingOrders []Order `json:"missing_orders"`
	}
	var orders []Order 	
	if err:=c.ShouldBindJSON(&orders);err!=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error binding body data",
			"error":   err.Error(),
		})
		return
	}
	preloadedDb=h.DB.Model(&models.Tenant{}).Preload("Sellers").Preload("Hubs").Preload("Skus")
	for _,o:=range(orders){
		
	}
}
