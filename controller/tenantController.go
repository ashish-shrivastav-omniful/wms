package controller

import (
	"context"
	"net/http"
	"wms/models"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type TenantHandler struct {
	DB  *gorm.DB
	RC  redis.Client
	Ctx *context.Context
}

func CreateTenantController(db *gorm.DB, rd redis.Client, ctx *context.Context) *TenantHandler {
	return &TenantHandler{DB: db, RC: rd, Ctx: ctx}
}

func (h *TenantHandler) ViewTenants(c *gin.Context) {
	var tenants []models.Tenant
	if results := h.DB.Find(&tenants); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during fetching hubs data",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, tenants)
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var tenantjson struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&tenantjson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during binding body data",
			"error":   err.Error(),
		})
		return
	}
	var newTenant models.Tenant
	newTenant.Name = tenantjson.Name
	if results := h.DB.Save(&newTenant); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error saving data in database",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "new Tenant created",
	})
}
