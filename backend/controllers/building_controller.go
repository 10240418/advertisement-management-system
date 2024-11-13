package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
)

// CreateBuildingInput 定义创建大厦的输入结构体
type CreateBuildingInput struct {
	Name             string `json:"name" binding:"required"`
	Address          string `json:"address"`
	AdvertisementIDs []uint `json:"advertisement_ids"`
	// NoticeIDs 如果需要关联通知，也可以添加
}

// CreateBuilding 创建新大厦，并关联广告
func CreateBuilding(c *gin.Context) {
	var input CreateBuildingInput

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 初始化大厦实例
	building := models.Building{
		Name:    input.Name,
		Address: input.Address,
	}

	// 保存大厦到数据库
	if err := config.DB.Create(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建大厦失败"})
		return
	}

	// 处理广告关联
	if len(input.AdvertisementIDs) > 0 {
		var ads []models.Advertisement
		if err := config.DB.Where("id IN ?", input.AdvertisementIDs).Find(&ads).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询广告失败"})
			return
		}

		if len(ads) != len(input.AdvertisementIDs) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "某些广告 ID 不存在"})
			return
		}

		// 创建 AdvertisementBuilding 关联记录，PlayDuration 默认值为广告的 VideoDuration
		for _, ad := range ads {
			association := models.AdvertisementBuilding{
				AdvertisementID: ad.ID,
				BuildingID:      building.ID,
				PlayDuration:    ad.VideoDuration, // 默认设置为广告的 VideoDuration
			}
			if err := config.DB.Create(&association).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "关联大厦与广告失败"})
				return
			}
		}
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings").Where("id = ?", building.ID).First(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		return
	}

	c.JSON(http.StatusCreated, building)
}

// UpdateBuildingInput 定义更新大厦的输入结构体
type UpdateBuildingInput struct {
	Name             string `json:"name"`
	Address          string `json:"address"`
	AdvertisementIDs []uint `json:"advertisement_ids"`
	// NoticeIDs 如果需要关联通知，也可以添加
}

// UpdateBuilding 更新大厦信息，并更新与广告的关联
func UpdateBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building

	// 查找大厦
	if err := config.DB.First(&building, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "大厦未找到"})
		return
	}

	var input UpdateBuildingInput

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if input.Name != "" {
		building.Name = input.Name
	}
	if input.Address != "" {
		building.Address = input.Address
	}

	// 保存更新后的大厦
	if err := config.DB.Save(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新大厦失败"})
		return
	}

	// 处理广告关联
	if input.AdvertisementIDs != nil {
		// 清除现有关联
		if err := config.DB.Where("building_id = ?", building.ID).Delete(&models.AdvertisementBuilding{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "清除现有关联失败"})
			return
		}

		// 查询新的广告
		var ads []models.Advertisement
		if err := config.DB.Where("id IN ?", input.AdvertisementIDs).Find(&ads).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询广告失败"})
			return
		}

		// 验证所有广告 ID 是否存在
		if len(ads) != len(input.AdvertisementIDs) {
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
			if err := config.DB.Create(&association).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "关联广告与大厦失败"})
				return
			}
		}
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings").Where("id = ?", building.ID).First(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		return
	}

	c.JSON(http.StatusOK, building)
}

// GetBuildings 获取所有大厦，并支持分页和排序
func GetBuildings(c *gin.Context) {
	var buildings []models.Building
	var count int64

	// 从查询参数中获取分页信息
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	desc := c.DefaultQuery("desc", "true")

	// 计算偏移量
	offset := (pageNum - 1) * pageSize

	// 构建查询
	query := config.DB.Model(&models.Building{})

	// 添加排序
	if strings.ToLower(desc) == "true" {
		query = query.Order("name DESC")
	} else {
		query = query.Order("name ASC")
	}

	// 执行查询并进行分页
	if err := query.Offset(offset).Limit(pageSize).Preload("AdvertisementBuildings").Preload("BuildingNotices").Find(&buildings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		return
	}

	// 获取总记录数用于分页
	query.Count(&count)

	// 返回数据和分页信息
	c.JSON(http.StatusOK, gin.H{
		"data":     buildings,
		"total":    count,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	})
}

// GetBuilding 获取单个大厦，并预加载关联的广告及播放时长
func GetBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building
	if err := config.DB.Preload("AdvertisementBuildings").Preload("BuildingNotices").First(&building, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "大厦未找到"})
		return
	}
	c.JSON(http.StatusOK, building)
}

// DeleteBuilding 删除大厦（硬删除）
func DeleteBuilding(c *gin.Context) {
	id := c.Param("id")

	// 硬删除大厦记录
	if err := config.DB.Unscoped().Delete(&models.Building{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除大厦失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "大厦删除成功"})
}
