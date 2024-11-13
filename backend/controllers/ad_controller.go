package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
)

// CreateAdInput 定义创建广告的输入结构体
type CreateAdInput struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	VideoURL      string `json:"video_url"`
	Status        string `json:"status" binding:"required"`
	VideoDuration int64  `json:"video_duration" binding:"required"` // 以秒为单位
	BuildingIDs   []uint `json:"building_ids" binding:"required"`
}

// CreateAd 创建广告，并建立与建筑的关联，设置 PlayDuration 为 VideoDuration
func CreateAd(c *gin.Context) {
	var input CreateAdInput

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询关联的建筑
	var buildings []models.Building
	if err := config.DB.Where("id IN ?", input.BuildingIDs).Find(&buildings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询建筑失败"})
		return
	}

	// 验证所有建筑 ID 是否存在
	if len(buildings) != len(input.BuildingIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "某些建筑 ID 不存在"})
		return
	}

	// 创建广告实例
	ad := models.Advertisement{
		Title:         input.Title,
		Description:   input.Description,
		ImageURL:      input.ImageURL,
		VideoURL:      input.VideoURL,
		Status:        input.Status,
		VideoDuration: input.VideoDuration,
	}

	// 保存广告到数据库
	if err := config.DB.Create(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建广告失败"})
		return
	}

	// 创建 AdvertisementBuilding 关联记录
	for _, building := range buildings {
		association := models.AdvertisementBuilding{
			AdvertisementID: ad.ID,
			BuildingID:      building.ID,
			PlayDuration:    ad.VideoDuration, // 默认为 VideoDuration
		}
		if err := config.DB.Create(&association).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "关联广告与建筑失败"})
			return
		}
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings").Where("id = ?", ad.ID).First(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		return
	}

	c.JSON(http.StatusCreated, ad)
}

// UpdateAdInput 定义更新广告的输入结构体
type UpdateAdInput struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	VideoURL      string `json:"video_url"`
	Status        string `json:"status"`
	VideoDuration int64  `json:"video_duration"` // 以秒为单位
	BuildingIDs   []uint `json:"building_ids"`
}

// UpdateAd 更新广告，并更新与建筑的关联
func UpdateAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Advertisement

	// 查找广告
	if err := config.DB.Preload("AdvertisementBuildings").First(&ad, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "广告未找到"})
		return
	}

	var input UpdateAdInput

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if input.Title != "" {
		ad.Title = input.Title
	}
	if input.Description != "" {
		ad.Description = input.Description
	}
	if input.ImageURL != "" {
		ad.ImageURL = input.ImageURL
	}
	if input.VideoURL != "" {
		ad.VideoURL = input.VideoURL
	}
	if input.Status != "" {
		ad.Status = input.Status
	}
	if input.VideoDuration != 0 {
		ad.VideoDuration = input.VideoDuration
	}

	// 保存更新后的广告
	if err := config.DB.Save(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新广告失败"})
		return
	}

	// 处理建筑关联
	if input.BuildingIDs != nil {
		// 清除现有关联
		if err := config.DB.Where("advertisement_id = ?", ad.ID).Delete(&models.AdvertisementBuilding{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "清除现有关联失败"})
			return
		}

		// 查询新的建筑
		var buildings []models.Building
		if err := config.DB.Where("id IN ?", input.BuildingIDs).Find(&buildings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询建筑失败"})
			return
		}

		// 验证所有建筑 ID 是否存在
		if len(buildings) != len(input.BuildingIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "某些建筑 ID 不存在"})
			return
		}

		// 创建新的关联记录
		for _, building := range buildings {
			association := models.AdvertisementBuilding{
				AdvertisementID: ad.ID,
				BuildingID:      building.ID,
				PlayDuration:    ad.VideoDuration, // 默认为 VideoDuration
			}
			if err := config.DB.Create(&association).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "关联广告与建筑失败"})
				return
			}
		}
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings").Where("id = ?", ad.ID).First(&ad).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		return
	}

	c.JSON(http.StatusOK, ad)
}

// GetAds 获取所有广告，并支持分页和排序
func GetAds(c *gin.Context) {
	var ads []models.Advertisement
	var count int64

	// 从查询参数中获取分页信息
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	desc := c.DefaultQuery("desc", "true")

	// 计算偏移量
	offset := (pageNum - 1) * pageSize

	// 构建查询
	query := config.DB.Model(&models.Advertisement{})

	// 添加排序
	if strings.ToLower(desc) == "true" {
		query = query.Order("username DESC")
	} else {
		query = query.Order("username ASC")
	}

	// 执行查询并进行分页
	if err := query.Offset(offset).Limit(pageSize).Preload("AdvertisementBuildings").Find(&ads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		return
	}

	// 获取总记录数用于分页
	query.Count(&count)

	// 返回数据和分页信息
	c.JSON(http.StatusOK, gin.H{
		"data":     ads,
		"total":    count,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	})
}

// DeleteAd 删除广告（硬删除）
func DeleteAd(c *gin.Context) {
	id := c.Param("id")

	// 硬删除广告记录
	if err := config.DB.Unscoped().Delete(&models.Advertisement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除广告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "广告删除成功"})
}

// UpdatePlayDuration 更新广告与建筑之间的播放时长
func UpdatePlayDuration(c *gin.Context) {
	var input struct {
		AdvertisementID uint  `json:"advertisement_id" binding:"required"`
		BuildingID      uint  `json:"building_id" binding:"required"`
		PlayDuration    int64 `json:"play_duration" binding:"required"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找关联记录
	var association models.AdvertisementBuilding
	if err := config.DB.Where("advertisement_id = ? AND building_id = ?", input.AdvertisementID, input.BuildingID).First(&association).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "关联记录未找到"})
		return
	}

	// 更新 PlayDuration
	association.PlayDuration = input.PlayDuration
	if err := config.DB.Save(&association).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新播放时长失败"})
		return
	}

	c.JSON(http.StatusOK, association)
}
