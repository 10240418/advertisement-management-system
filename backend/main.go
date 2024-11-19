package main

import (
	"log"

	"time"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/models"
	"github.com/10240418/advertisement-management-system/backend/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatalf("加载环境变量失败: %v", err)
	}

	// 初始化数据库连接
	if err := config.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 运行数据库迁移
	if err := config.DB.AutoMigrate(&models.Advertisement{}, &models.Building{}, &models.AdvertisementBuilding{}, &models.Administrator{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置路由
	r := routers.SetupRouter()

	// 配置 CORS
	configCORSMiddleware(r)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// configCORSMiddleware 配置 CORS 中间件
func configCORSMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 根据需求调整允许的源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
