package main

import (
	"log"
	"os"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/controllers"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/gin-contrib/cors" // 导入CORS中间件
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 自动迁移
	config.DB.AutoMigrate(&models.Advertisement{})

	// 设置Gin为release模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 配置CORS，允许任何来源
	configCors := cors.Config{
		AllowAllOrigins: true, // 允许所有来源
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}

	r.Use(cors.New(configCors)) // 应用CORS中间件

	// 定义API路由
	api := r.Group("/api")
	{
		api.GET("/ads", controllers.GetAds)
		api.GET("/ads/:id", controllers.GetAd)
		api.POST("/ads", controllers.CreateAd)
		api.PUT("/ads/:id", controllers.UpdateAd)
		api.DELETE("/ads/:id", controllers.DeleteAd)
	}

	// 读取环境变量PORT，默认为8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
