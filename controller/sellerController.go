package controller

import (
	"context"
	"net/http"
	"wms/models"

	"github.com/gin-gonic/gin"
	// "github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type SellerHandler struct {
	DB  *gorm.DB
	// RC  redis.Client
	Ctx *context.Context
}

func CreateSellerController(db *gorm.DB, ctx *context.Context) *SellerHandler {
	return &SellerHandler{DB: db, Ctx: ctx}
}

func (h *SellerHandler) ViewSellers(c *gin.Context) {
	var sellers []models.Seller
	if results := h.DB.Find(&sellers); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during fetching hubs data",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, sellers)
}

func (h *SellerHandler) CreateSeller(c *gin.Context) {
	var sellerjson struct {
		TenantId uint   `json:"tenantid" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Phone    string `json:"phone" binding:"required"`
		Address  string `json:"address" binding:"required"`
	}
	if err := c.ShouldBindJSON(&sellerjson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error during binding body data",
			"error":   err.Error(),
		})
		return
	}
	var newSeller models.Seller
	newSeller.Address = sellerjson.Address
	newSeller.Name = sellerjson.Name
	newSeller.TenantId = sellerjson.TenantId
	newSeller.Email = sellerjson.Email
	newSeller.Phone = sellerjson.Phone
	if results := h.DB.Save(&newSeller); results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error saving data in database",
			"error":   results.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "new seller created",
	})
}
