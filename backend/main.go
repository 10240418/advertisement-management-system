package main

import (
	"log"

	"github.com/10240418/advertisement-management-system/backend/config"
	"github.com/10240418/advertisement-management-system/backend/controllers"
	"github.com/10240418/advertisement-management-system/backend/middleware"
	"github.com/gin-contrib/cors" // 导入CORS中间件
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config.InitDB()

	r := gin.Default()
	// 配置CORS，允许任何来源
	configCors := cors.Config{
		AllowAllOrigins: true, // 允许所有来源
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}

	r.Use(cors.New(configCors)) // 应用CORS中间件

	// 公共路由
	r.POST("/admin/register", controllers.RegisterAdmin)
	r.POST("/admin/login", controllers.LoginAdmin)
	r.GET("/admin/users", controllers.GetAdminUsers)
	r.DELETE("/admin/users", controllers.DeleteAdmin)
	r.PUT("/admin/user", controllers.UpdateAdminPassword)

	// 受保护的路由组
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware()) // 应用认证中间件

	{
		// 广告路由
		protected.GET("/ads", controllers.GetAds)
		protected.GET("/ads/:id", controllers.GetAd)
		protected.POST("/ads", controllers.CreateAd)
		protected.PUT("/ads/:id", controllers.UpdateAd)
		protected.DELETE("/ads/:id", controllers.DeleteAd)

		// 大厦路由
		protected.GET("/buildings", controllers.GetBuildings)
		protected.GET("/buildings/:id", controllers.GetBuilding)
		protected.POST("/buildings", controllers.CreateBuilding)
		protected.PUT("/buildings/:id", controllers.UpdateBuilding)
		protected.DELETE("/buildings/:id", controllers.DeleteBuilding)

		// 通知路由
		protected.GET("/notices", controllers.GetNotices)
		protected.GET("/notices/:id", controllers.GetNotice)
		protected.POST("/notices", controllers.CreateNotice)
		protected.PUT("/notices/:id", controllers.UpdateNotice)
		protected.DELETE("/notices/:id", controllers.DeleteNotice)
	}

	r.Run(":8080") // 监听端口
}
