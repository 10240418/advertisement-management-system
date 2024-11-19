package controllers

import (
	"net/http"
	"strconv"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateAdInput 定义更新广告的输入结构体
type UpdateAdInput struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	VideoURL      string `json:"video_url"`
	Status        string `json:"status"`
	VideoDuration int64  `json:"video_duration"` // 以秒为单位
}

// UpdatePlayDurationInput 定义更新播放时长的输入结构体
type UpdatePlayDurationInput struct {
	AdvertisementID uint  `json:"advertisement_id" binding:"required"`
	BuildingID      uint  `json:"building_id" binding:"required"`
	PlayDuration    int64 `json:"play_duration" binding:"required"`
}

// GetAds 获取所有广告，并支持分页和排序
func GetAds(c *gin.Context) {
	// 从查询参数中获取分页信息
	pageNumStr := c.DefaultQuery("pageNum", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	desc := c.DefaultQuery("desc", "true")

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		pageNum = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	offset := (pageNum - 1) * pageSize

	var ads []models.Advertisement
	var count int64

	// 构建基础查询
	baseQuery := config.DB.Model(&models.Advertisement{})

	// 添加排序
	if desc == "true" {
		baseQuery = baseQuery.Order("created_at DESC")
	} else {
		baseQuery = baseQuery.Order("created_at ASC")
	}

	// 获取总记录数不包含limit和offset
	if err := baseQuery.Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取总数失败"})
		return
	}

	// 执行查询并进行分页
	if err := baseQuery.Offset(offset).Limit(pageSize).Find(&ads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		return
	}

	// 返回数据和分页信息
	c.JSON(http.StatusOK, gin.H{
		"data":     ads,
		"total":    count,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	})
}

// GetAd 获取单个广告
func GetAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Advertisement

	// 查找广告并预加载关联的 AdvertisementBuildings 和 Building
	if err := config.DB.
		Preload("AdvertisementBuildings.Building").
		First(&ad, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		}
		return
	}

	c.JSON(http.StatusOK, ad)
}

// CreateAd 创建新广告
func CreateAd(c *gin.Context) {
	var input models.Advertisement

	// 绑定 JSON 数据到 Advertisement 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 创建广告
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建广告失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings.Building").First(&input, input.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// UpdateAd 更新广告的基本信息
func UpdateAd(c *gin.Context) {
	id := c.Param("id")
	var ad models.Advertisement

	// 查找广告
	if err := config.DB.First(&ad, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "广告未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		}
		return
	}

	var input UpdateAdInput

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新广告字段
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

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 保存更新后的广告
	if err := tx.Save(&ad).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新广告失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings.Building").First(&ad, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取广告失败"})
		return
	}

	c.JSON(http.StatusOK, ad)
}

// DeleteAd 删除广告
func DeleteAd(c *gin.Context) {
	id := c.Param("id")

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 删除关联的 AdvertisementBuilding 记录
	if err := tx.Where("advertisement_id = ?", id).Delete(&models.AdvertisementBuilding{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除关联广告与建筑失败"})
		return
	}

	// 删除广告记录
	if err := tx.Unscoped().Delete(&models.Advertisement{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除广告失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "广告删除成功"})
}

// AddAdsToBuilding 通过 Building ID 添加多个 Advertisement 关联
func AddAdsToBuilding(c *gin.Context) {
	buildingID := c.Param("id")
	var input struct {
		AdvertisementIDs []uint `json:"advertisement_ids" binding:"required"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 查找建筑
	var building models.Building
	if err := tx.First(&building, buildingID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "建筑未找到"})
		return
	}

	// 查找广告
	var ads []models.Advertisement
	if err := tx.Where("id IN ?", input.AdvertisementIDs).Find(&ads).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询广告失败"})
		return
	}

	// 验证所有广告 ID 是否存在
	if len(ads) != len(input.AdvertisementIDs) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "某些广告 ID 不存在"})
		return
	}

	// 创建新的关联记录
	for _, ad := range ads {
		association := models.AdvertisementBuilding{
			AdvertisementID: ad.ID,
			BuildingID:      building.ID,
			PlayDuration:    ad.VideoDuration, // 默认为 VideoDuration
		}
		if err := tx.Create(&association).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "关联广告与建筑失败"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "广告与建筑关联成功"})
}

// RemoveAdsFromBuilding 通过 Building ID 删除多个 Advertisement 关联
func RemoveAdsFromBuilding(c *gin.Context) {
	buildingID := c.Param("id")
	var input struct {
		AdvertisementIDs []uint `json:"advertisement_ids" binding:"required"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 查找建筑
	var building models.Building
	if err := tx.First(&building, buildingID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "建筑未找到"})
		return
	}

	// 删除指定的关联记录
	if err := tx.Where("building_id = ? AND advertisement_id IN ?", building.ID, input.AdvertisementIDs).Delete(&models.AdvertisementBuilding{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除关联失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "广告与建筑关联删除成功"})
}

// AddBuildingsToAd 通过 Advertisement ID 添加多个 Building 关联
func AddBuildingsToAd(c *gin.Context) {
	adID := c.Param("id")
	var input struct {
		BuildingIDs []uint `json:"building_ids" binding:"required"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 查找广告
	var ad models.Advertisement
	if err := tx.First(&ad, adID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "广告未找到"})
		return
	}

	// 查找建筑
	var buildings []models.Building
	if err := tx.Where("id IN ?", input.BuildingIDs).Find(&buildings).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询建筑失败"})
		return
	}

	// 验证所有建筑 ID 是否存在
	if len(buildings) != len(input.BuildingIDs) {
		tx.Rollback()
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
		if err := tx.Create(&association).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "关联广告与建筑失败"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "广告与建筑关联成功"})
}

// RemoveBuildingsFromAd 通过 Advertisement ID 删除多个 Building 关联
func RemoveBuildingsFromAd(c *gin.Context) {
	adID := c.Param("id")
	var input struct {
		BuildingIDs []uint `json:"building_ids" binding:"required"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 查找广告
	var ad models.Advertisement
	if err := tx.First(&ad, adID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "广告未找到"})
		return
	}

	// 删除指定的关联记录
	if err := tx.Where("advertisement_id = ? AND building_id IN ?", ad.ID, input.BuildingIDs).Delete(&models.AdvertisementBuilding{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除关联失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "广告与建筑关联删除成功"})
}

// GetAdvertisementsByBuilding 获取指定 Building ID 关联的所有 Advertisement 对象
func GetAdvertisementsByBuilding(c *gin.Context) {
	buildingID := c.Param("id")
	var associations []models.AdvertisementBuilding

	// 查询关联记录
	if err := config.DB.Where("building_id = ?", buildingID).Find(&associations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询关联失败"})
		return
	}

	// 提取 Advertisement IDs
	var adIDs []uint
	for _, assoc := range associations {
		adIDs = append(adIDs, assoc.AdvertisementID)
	}

	// 查询 Advertisement 对象
	var advertisements []models.Advertisement
	if err := config.DB.Where("id IN ?", adIDs).Find(&advertisements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询广告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"advertisements": advertisements,
	})
}

// GetBuildingsByAdvertisement 获取指定 Advertisement ID 关联的所有 Building 对象
func GetBuildingsByAdvertisement(c *gin.Context) {
	adID := c.Param("id")
	var associations []models.AdvertisementBuilding

	// 查询关联记录
	if err := config.DB.Where("advertisement_id = ?", adID).Find(&associations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询关联失败"})
		return
	}

	// 提取 Building IDs
	var buildingIDs []uint
	for _, assoc := range associations {
		buildingIDs = append(buildingIDs, assoc.BuildingID)
	}

	// 查询 Building 对象
	var buildings []models.Building
	if err := config.DB.Where("id IN ?", buildingIDs).Find(&buildings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询建筑失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"buildings": buildings,
	})
}

// UpdatePlayDuration 更新广告与建筑之间的播放时长
func UpdatePlayDuration(c *gin.Context) {
	var input UpdatePlayDurationInput

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
