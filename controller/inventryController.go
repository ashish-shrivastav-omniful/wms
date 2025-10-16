package controller

import (
	"context"
	"net/http"
	"wms/models"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type InventoryHandler struct{
	DB *gorm.DB
	RC redis.Client
	Ctx *context.Context
}

func CreateInventoryController (db *gorm.DB,rd redis.Client,ctx *context.Context)*InventoryHandler{
	return &InventoryHandler{DB: db,RC: rd,Ctx: ctx}
}
//to get inventory details for each sku at each hub
func (h *InventoryHandler) viewInventory (c *gin.Context){
	var inventories [] models.Inventory
	if results:=h.DB.Find(&inventories); results.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Error during fetching hubs data",
			"error":results.Error.Error(),
		})
		return 
	}
	c.JSON(200,inventories)


}
//to change inventory details for a particular sku
func (h *InventoryHandler)editInventory (c *gin.Context){
	var inventoryjson struct{
		HubId uint `json:"hubid" binding:"required"`
		SkuId uint `json:"skuid" binding:"required"`
		Qty uint `json:"qty" binding:"required"`
	}
	if err:=c.ShouldBindJSON(&inventoryjson);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Error binding sku data",
			"error":err.Error(),
		})
		return 
	}
	var editedInventory models.Inventory
	if err:=h.DB.Where("skuid = ? AND hubid = ?",inventoryjson.SkuId,inventoryjson.HubId).First(&editedInventory);err.Error!=nil{
		editedInventory.HubId=inventoryjson.HubId
		editedInventory.SkuId=inventoryjson.SkuId
		editedInventory.Qty=inventoryjson.Qty
		if err:=h.DB.Save(&editedInventory);err!=nil {
			c.JSON(http.StatusBadRequest,gin.H{
			"message":"Could not find inventry and error creating new inventory",
			"error":err.Error.Error(),
		})
		return
		}
		c.JSON(200,gin.H{
			"message":"Could not find inventory data so created new entry",
		})
		return 
	}
	editedInventory.Qty=inventoryjson.Qty
	if err:=h.DB.Save(&editedInventory); err.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Error saving edited data",
			"error":err.Error.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"message":"Edited inventry details",
	})
}
//to update inventory data after sales of particular sku
func (h *InventoryHandler)updateInventoryAftersales(c *gin.Context){
	var inventoryjson struct{
		HubId uint `json:"hubid" binding:"required"`
		SkuId uint `json:"skuid" binding:"required"`
		Qty uint `json:"qty" binding:"required"`
	}
	if err:=c.ShouldBindJSON(&inventoryjson);err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Error binding item data",
			"error":err.Error(),
		})
		return 
	}
	var editedInventory models.Inventory
	if err:=h.DB.Where("skuid = ? AND hubid = ?",inventoryjson.SkuId,inventoryjson.HubId).First(&editedInventory);err.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Error finding inventory data",
			"error":err.Error.Error(),
		})
		return 
	}
	if(editedInventory.Qty<inventoryjson.Qty){
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Sufficient items are not there in the inventory",
		})
		return
	}
	editedInventory.Qty-=inventoryjson.Qty
	if err:=h.DB.Save(&editedInventory); err.Error!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Error saving edited data",
			"error":err.Error.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"message":"Updated inventry data after sales",
	})
}