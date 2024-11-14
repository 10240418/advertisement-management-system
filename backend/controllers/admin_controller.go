package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAdmin 注册新的管理员
func RegisterAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 修剪用户名和密码，去除前后空格
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// 检查用户名和密码为空
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	// 检查用户名是否已存在
	var existingAdmin models.Administrator
	if err := config.DB.Where("username = ?", input.Username).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Password encryption failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	admin := models.Administrator{
		Username: input.Username,
		Password: string(hashedPassword), // 确保这里存储的是生成的哈希值
	}

	// 保存管理员到数据库
	if err := config.DB.Create(&admin).Error; err != nil {
		fmt.Printf("Failed to create administrator: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册管理员失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "管理员注册成功"})
}

// LoginAdmin 管理员登录
func LoginAdmin(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 修剪用户名和密码，去除前后空格
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// 检查用户名和密码是否为空
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	// 查询数据库中的管理员
	var admin models.Administrator
	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	// 生成JWT
	token, err := config.GenerateToken(admin.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token})
}

// GetAdminUsers 获取所有管理员
func GetAdminUsers(c *gin.Context) {
	// 从查询参数中获取分页信息
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	desc := c.DefaultQuery("desc", "true")

	// 计算分页的偏移量
	offset := (pageNum - 1) * pageSize

	var admins []models.Administrator
	var count int64

	// 构建查询
	query := config.DB.Model(&models.Administrator{})

	// 添加排序
	if desc == "true" {
		query = query.Order("username DESC")
	} else {
		query = query.Order("username ASC")
	}

	// 执行查询并进行分页
	if err := query.Offset(offset).Limit(pageSize).Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取管理员失败"})
		return
	}

	// 获取总记录数用于分页
	query.Count(&count)

	// 返回数据和分页信息
	c.JSON(http.StatusOK, gin.H{
		"data":     admins,
		"total":    count,
		"pageNum":  pageNum,
		"pageSize": pageSize,
	})
}

// 修改密码
func UpdateAdminPassword(c *gin.Context) {
	var input struct {
		Username    string `json:"username"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查询数据库中的管理员
	var admin models.Administrator
	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码不正确"})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "新密码加密失败"})
		return
	}

	// 更新密码
	admin.Password = string(hashedPassword)
	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
}

// DeleteAdmin 删除管理员（硬删除）
func DeleteAdmin(c *gin.Context) {
	var input struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 是必需的"})
		return
	}

	// 硬删除管理员记录
	if err := config.DB.Unscoped().Delete(&models.Administrator{}, input.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除管理员失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "管理员删除成功"})
}
