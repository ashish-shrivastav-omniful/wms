package controller

import (
	"context"
	"net/http"
	"wms/models"

	"github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type HubController struct {
	DB  *gorm.DB
	// RC  redis.Client
	Ctx *context.Context
}

func CreateHubController(db *gorm.DB, ctx *context.Context) *HubController {
	return &HubController{DB: db, Ctx: ctx}
}

// getting all hubs
func (h *HubController) ViewHubs(c *gin.Context) {
	var hubs []models.Hub
	if results := h.DB.Find(&hubs); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during fetching hubs data",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, hubs)
}

// creating a hub
func (h *HubController) CreateHub(c *gin.Context) {
	var hubjson struct {
		TenantId uint   `json:"tenantid" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Address  string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&hubjson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during binding body data",
			"error":   err.Error(),
		})
		return
	}
	var newhub models.Hub
	newhub.Address = hubjson.Address
	newhub.Name = hubjson.Name
	newhub.TenantId = hubjson.TenantId
	if results := h.DB.Save(&newhub); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error saving data in database",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "new Hub created",
	})
}
