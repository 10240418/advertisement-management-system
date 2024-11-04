package controllers

import (
	"net/http"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
)

func GetAds(c *gin.Context) {
	var ads []models.Advertisement
	config.DB.Find(&ads)
	c.JSON(http.StatusOK, ads)
}

func GetAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Advertisement
	if err := config.DB.First(&ad, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
		return
	}
	c.JSON(http.StatusOK, ad)
}

func CreateAd(c *gin.Context) {
	var ad models.Advertisement
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&ad)
	c.JSON(http.StatusCreated, ad)
}

func UpdateAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Advertisement
	if err := config.DB.First(&ad, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
		return
	}
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&ad)
	c.JSON(http.StatusOK, ad)
}

func DeleteAd(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Advertisement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ad deleted"})
}
