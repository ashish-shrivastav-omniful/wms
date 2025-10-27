package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"wms/models"

	"github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type InventoryHandler struct {
	DB  *gorm.DB
	// RC  redis.Client
	Ctx *context.Context
}

func CreateInventoryController(db *gorm.DB, ctx *context.Context) *InventoryHandler {
	return &InventoryHandler{DB: db, Ctx: ctx}
}

// to get inventory details for each sku at each hub
func (h *InventoryHandler) ViewInventory(c *gin.Context) {
	var inventories []models.Inventory
	if results := h.DB.Find(&inventories); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during fetching hubs data",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, inventories)

}

// to change inventory details for a particular sku
func (h *InventoryHandler) EditInventory(c *gin.Context) {
	var input struct {
		HubId uint `json:"hubid" binding:"required"`
		SkuId uint `json:"skuid" binding:"required"`
		Qty   uint `json:"qty" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error binding SKU data",
			"error":   err.Error(),
		})
		return
	}

	var inventory models.Inventory
	result := h.DB.Where("sku_id = ? AND hub_id = ?", input.SkuId, input.HubId).First(&inventory)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create new inventory if not found
		inventory = models.Inventory{
			HubId: input.HubId,
			SkuId: input.SkuId,
			Qty:   input.Qty,
		}
		if err := h.DB.Create(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Could not create inventory",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Inventory not found, created new entry",
			"data":    inventory,
		})
		return
	} else if result.Error != nil {
		// Any other DB error
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
			"error":   result.Error.Error(),
		})
		return
	}

	// Record exists, update quantity
	inventory.Qty = input.Qty
	if err := h.DB.Save(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error updating inventory",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory updated successfully",
		"data":    inventory,
	})
}


// to update inventory data after sales of particular sku
func (h *InventoryHandler) UpdateInventoryAftersales(c *gin.Context) {
	var inventoryjson struct {
		OrderId string `json:"order_id"`
		HubId string `json:"hubid"`
		SkuId string `json:"item_id" binding:"required"`
		Qty   string `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&inventoryjson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error binding item data",
			"error":   err.Error(),
		})
		return
	}
	var editedInventory models.Inventory
	
	if err := h.DB.Joins("Join skus ON inventories.sku_id = skus.id").Where("sku_code = ? ", inventoryjson.SkuId).First(&editedInventory); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error finding inventory data",
			"error":   err.Error.Error(),
		})
		return
	}
	q,_:=strconv.ParseUint(inventoryjson.Qty,10,64)
	if editedInventory.Qty <  uint(q){
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Sufficient items are not there in the inventory",
		})
		return
	}
	editedInventory.Qty -= uint(q)
	if err := h.DB.Save(&editedInventory); err.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error saving edited data",
			"error":   err.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Updated inventry data after sales",
	})
}
