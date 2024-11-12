package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAdmin 注册新的管理员
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
	// 查询数据库中的管理员
	var admin models.Administrator

	if err := config.DB.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		fmt.Printf("Login failed: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	fmt.Printf("Fetched admin from DB: %+v\n", admin)

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		fmt.Printf("Password mismatch for user '%s': %v\n", input.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		fmt.Printf("Password mismatch for user '%s': %v\n", []byte(admin.Password), err)
		fmt.Printf("Password mismatch for user '%s': %v\n", []byte(input.Password), err)
		var password = "healthist"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Error generating hash:", err)
			return
		}
		fmt.Println("Generated hash:", string(hashedPassword))
		return
	} else {
		fmt.Printf("Password match for user '%s'\n", input.Username)
	}

	// 生成JWT
	token, err := config.GenerateToken(admin.Username)
	if err != nil {
		fmt.Printf("Token generation failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	fmt.Printf("Token generated for user '%s': %s\n", input.Username, token)

	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "token": token})
}
