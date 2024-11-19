package routers

import (
	"github.com/10240418/advertisement-management-system/backend/controllers"
	"github.com/10240418/advertisement-management-system/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 公共路由
	r.POST("/admin/register", controllers.RegisterAdmin)
	r.POST("/admin/login", controllers.LoginAdmin)

	// 受保护的路由组
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware()) // 应用认证中间件
	{
		// 广告路由
		ads := protected.Group("/ads")
		{
			ads.GET("", controllers.GetAds)
			ads.GET("/:id", controllers.GetAd)
			ads.POST("", controllers.CreateAd)
			ads.PUT("/:id", controllers.UpdateAd)
			ads.DELETE("/:id", controllers.DeleteAd)

			// 新增的路由：管理广告与建筑的关联
			ads.POST("/:id/buildings", controllers.AddBuildingsToAd)           // 添加建筑到广告
			ads.DELETE("/:id/buildings", controllers.RemoveBuildingsFromAd)    // 删除建筑与广告的关联
			ads.GET("/:id/buildings", controllers.GetBuildingsByAdvertisement) // 获取广告关联的建筑 IDs
		}

		// 大厦路由
		buildings := protected.Group("/buildings")
		{
			buildings.GET("", controllers.GetBuildings)
			buildings.GET("/:id", controllers.GetBuilding)
			buildings.POST("", controllers.CreateBuilding)
			buildings.PUT("/:id", controllers.UpdateBuilding)
			buildings.DELETE("/:id", controllers.DeleteBuilding)

			// 新增的路由：管理建筑与广告的关联
			buildings.POST("/:id/ads", controllers.AddAdsToBuilding)           // 添加广告到建筑
			buildings.DELETE("/:id/ads", controllers.RemoveAdsFromBuilding)    // 删除广告与建筑的关联
			buildings.GET("/:id/ads", controllers.GetAdvertisementsByBuilding) // 获取建筑关联的广告 IDs
		}

		// 管理员路由
		admins := protected.Group("/admins")
		{
			admins.GET("/users", controllers.GetAdminUsers)
			admins.DELETE("/users", controllers.DeleteAdmin)
			admins.PUT("/user", controllers.UpdateAdminPassword)
		}
	}

	return r
}
