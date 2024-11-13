package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
)

// GetNotices 获取所有通知，并支持分页和排序
func GetNotices(c *gin.Context) {
	var notices []models.Notice
	var count int64

	// 从查询参数中获取分页信息
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	desc := c.DefaultQuery("desc", "true")

	// 计算偏移量
	offset := (pageNum - 1) * pageSize

	// 构建查询
	query := config.DB.Model(&models.Notice{})

	// 添加排序
	if strings.ToLower(desc) == "true" {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("created_at ASC")
	}

	// 执行查询并进行分页
	if err := query.Preload("Buildings").Offset(offset).Limit(pageSize).Find(&notices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取通知失败"})
		return
	}

	// 获取总记录数用于分页
	query.Count(&count)

	// 返回数据和分页信息
	c.JSON(http.StatusOK, gin.H{
		"data":     notices,
		"total":    count,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	})
}

// GetNotice 获取单个通知
func GetNotice(c *gin.Context) {
	id := c.Param("id")
	var notice models.Notice
	if err := config.DB.Preload("Buildings").First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知未找到"})
		return
	}
	c.JSON(http.StatusOK, notice)
}

// CreateNotice 创建新通知
func CreateNotice(c *gin.Context) {
	var input struct {
		Title       string `json:"title" binding:"required"`
		PDFURL      string `json:"pdf_url"`
		BuildingIDs []uint `json:"building_ids" binding:"required"`
	}

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

	// 创建通知实例并关联建筑
	notice := models.Notice{
		Title:     input.Title,
		PDFURL:    input.PDFURL,
		Buildings: buildings,
	}

	// 保存通知到数据库
	if err := config.DB.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败"})
		return
	}

	c.JSON(http.StatusCreated, notice)
}

// UpdateNotice 更新通知
func UpdateNotice(c *gin.Context) {
	id := c.Param("id")
	var notice models.Notice
	if err := config.DB.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知未找到"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		PDFURL      string `json:"pdf_url"`
		BuildingIDs []uint `json:"building_ids"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if input.Title != "" {
		notice.Title = input.Title
	}
	if input.PDFURL != "" {
		notice.PDFURL = input.PDFURL
	}

	// 处理建筑关联
	if input.BuildingIDs != nil {
		// 清除现有关联
		if err := config.DB.Model(&notice).Association("Buildings").Clear(); err != nil {
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

		// 重新关联建筑
		if err := config.DB.Model(&notice).Association("Buildings").Append(buildings); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "关联建筑失败"})
			return
		}
	}

	// 保存更新后的通知
	if err := config.DB.Save(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新通知失败"})
		return
	}

	// 返回更新后的通知
	c.JSON(http.StatusOK, notice)
}

// DeleteNotice 删除通知
func DeleteNotice(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Notice{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "通知删除成功"})
}
