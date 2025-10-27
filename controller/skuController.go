package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"wms/models"

	"github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type SkuHandler struct {
	DB  *gorm.DB
	// RC  redis.Client
	Ctx *context.Context
}

func CreatSkuController(db *gorm.DB, ctx *context.Context) *SkuHandler {
	return &SkuHandler{DB: db, Ctx: ctx}
}

// view all skus present
func (h *SkuHandler) ViewSkus(c *gin.Context) {
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
func (h *SkuHandler) CreateSkus(c *gin.Context) {
	var skujson struct {
		TenantId uint   `json:"tenantid" binding:"required"`
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
	if err := h.DB.Save(&newSku); err.Error != nil {
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
func (h *SkuHandler) VerifySkus(c *gin.Context) {
	fmt.Println(c.Request.Body)
	type Order struct {
		SNo      string `json:"sno"`
		SellerID string   `json:"seller_id"`
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
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error binding body data",
			"error":   err.Error(),
		})
		return
	}
	var newOrderResponse OrderResponse
	preloadedDb := h.DB.Model(&models.Tenant{}).
		Joins("JOIN sellers ON sellers.tenant_id = tenants.id").
		Joins("JOIN skus ON skus.tenant_id = tenants.id").
		Joins("Join hubs ON hubs.tenant_id = tenants.id").
		Joins("Join inventories ON inventories.hub_id = hubs.id AND inventories.sku_id = skus.id")
	for _, o := range orders {
		results := preloadedDb.Where("Sellers.ID = ? and skus.Sku_Code = ?", o.SellerID, o.ItemID)
		if results.Error != nil {
			log.Println(results.Error.Error())
		} else {
			var cnt int64
			results.Count(&cnt)
			if cnt > 0 {
				newOrderResponse.ValidOrders = append(newOrderResponse.ValidOrders, o)
			} else {
				newOrderResponse.MissingOrders = append(newOrderResponse.MissingOrders, o)
			}
		}
		log.Println(o)
	}
	c.JSON(200, newOrderResponse)
}
