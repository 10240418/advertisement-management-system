package controllers

import (
	"net/http"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
)

// GetBuildings 获取所有大厦
func GetBuildings(c *gin.Context) {
	var buildings []models.Building
	config.DB.Preload("Advertisements").Preload("Notices").Find(&buildings)
	c.JSON(http.StatusOK, buildings)
}

// GetBuilding 获取单个大厦
func GetBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building
	if err := config.DB.Preload("Advertisements").Preload("Notices").First(&building, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
		return
	}
	c.JSON(http.StatusOK, building)
}

// CreateBuilding 创建新大厦
func CreateBuilding(c *gin.Context) {
	var building models.Building
	if err := c.ShouldBindJSON(&building); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create building"})
		return
	}
	c.JSON(http.StatusCreated, building)
}

// UpdateBuilding 更新大厦信息
func UpdateBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building
	if err := config.DB.First(&building, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
		return
	}
	if err := c.ShouldBindJSON(&building); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&building)
	c.JSON(http.StatusOK, building)
}

// DeleteBuilding 删除大厦
func DeleteBuilding(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Building{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Building deleted"})
}
