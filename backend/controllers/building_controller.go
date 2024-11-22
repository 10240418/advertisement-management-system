package controllers

import (
	"net/http"
	"strconv"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateBuildingInput 定义创建大厦的输入结构体
type CreateBuildingInput struct {
	Name             string `json:"name" binding:"required"`
	Address          string `json:"address"`
	BuildingID       string `json:"blg_id"`
	AdvertisementIDs []uint `json:"advertisement_ids"`
	// NoticeIDs 如果需要关联通知，也可以添加
}

// UpdateBuildingInput 定义更新大厦的输入结构体
type UpdateBuildingInput struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	BuildingID string `json:"blg_id"`
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
		Name:       input.Name,
		Address:    input.Address,
		BuildingID: input.BuildingID,
	}

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 保存大厦到数据库
	if err := tx.Create(&building).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建大厦失败"})
		return
	}

	// 处理广告关联
	if len(input.AdvertisementIDs) > 0 {
		var ads []models.Advertisement
		if err := tx.Where("id IN ?", input.AdvertisementIDs).Find(&ads).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询广告失败"})
			return
		}

		if len(ads) != len(input.AdvertisementIDs) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "某些广告 ID 不存在"})
			return
		}

		// 创建 AdvertisementBuilding 关联记录
		for _, ad := range ads {
			association := models.AdvertisementBuilding{
				AdvertisementID: ad.ID,
				BuildingID:      building.ID,
				PlayDuration:    ad.VideoDuration, // 默认为 VideoDuration
			}
			if err := tx.Create(&association).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "关联大厦与广告失败"})
				return
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings").First(&building, building.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		return
	}

	c.JSON(http.StatusOK, building)
}

// UpdateBuilding 更新大厦的基本信息
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

	// 更新大厦字段
	if input.Name != "" {
		building.Name = input.Name
	}
	if input.Address != "" {
		building.Address = input.Address
	}
	if input.BuildingID != "" {
		building.BuildingID = input.BuildingID
	}

	// 保存更新后的大厦
	if err := config.DB.Save(&building).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新大厦失败"})
		return
	}

	// 预加载关联数据返回
	if err := config.DB.Preload("AdvertisementBuildings").First(&building, building.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		return
	}

	c.JSON(http.StatusOK, building)
}

// DeleteBuilding 删除大厦
func DeleteBuilding(c *gin.Context) {
	id := c.Param("id")

	// 开始事务
	tx := config.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "启动事务失败"})
		return
	}

	// 删除关联的 AdvertisementBuilding 记录
	if err := tx.Where("building_id = ?", id).Delete(&models.AdvertisementBuilding{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除关联大厦与广告失败"})
		return
	}

	// 删除大厦记录
	if err := tx.Unscoped().Delete(&models.Building{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除大厦失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "大厦删除成功"})
}

// GetBuildings 获取所有大厦，并支持分页和排序
func GetBuildings(c *gin.Context) {
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

	var buildings []models.Building
	var count int64

	// 构建查询
	baseQuery := config.DB.Model(&models.Building{})

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
	if err := baseQuery.Offset(offset).Limit(pageSize).Find(&buildings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		return
	}

	// 返回数据和分页信息
	c.JSON(http.StatusOK, gin.H{
		"data":     buildings,
		"total":    count,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	})
}

// GetBuilding 获取单个大厦
func GetBuilding(c *gin.Context) {
	id := c.Param("id")
	var building models.Building

	if err := config.DB.
		Preload("AdvertisementBuildings").
		Preload("AdvertisementBuildings.Advertisement").
		Preload("AdvertisementBuildings.Building").
		First(&building, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "大厦未找到"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取大厦失败"})
		}
		return
	}

	c.JSON(http.StatusOK, building)
}
