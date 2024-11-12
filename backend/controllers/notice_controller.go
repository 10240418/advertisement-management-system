package controllers

import (
	"net/http"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
)

// GetNotices 获取所有通知
func GetNotices(c *gin.Context) {
	var notices []models.Notice
	config.DB.Preload("Building").Find(&notices)
	c.JSON(http.StatusOK, notices)
}

// GetNotice 获取单个通知
func GetNotice(c *gin.Context) {
	id := c.Param("id")
	var notice models.Notice
	if err := config.DB.Preload("Building").First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notice not found"})
		return
	}
	c.JSON(http.StatusOK, notice)
}

// CreateNotice 创建新通知
func CreateNotice(c *gin.Context) {
	var notice models.Notice
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notice"})
		return
	}
	c.JSON(http.StatusCreated, notice)
}

// UpdateNotice 更新通知
func UpdateNotice(c *gin.Context) {
	id := c.Param("id")
	var notice models.Notice
	if err := config.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notice not found"})
		return
	}
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&notice)
	c.JSON(http.StatusOK, notice)
}

// DeleteNotice 删除通知
func DeleteNotice(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Notice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notice deleted"})
}
