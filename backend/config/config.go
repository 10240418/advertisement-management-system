package config

import (
	"log"
	"os"

	"github.com/10240418/advertisement-management-system/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	dsn := os.Getenv("POSTGRES_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 确保表名为复数形式
			SingularTable: false,
		},
		Logger: logger.Default.LogMode(logger.Info), // 设置日志级别为 Info 以便调试
	})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	// 自动迁移模型，包括关联表
	err = DB.AutoMigrate(
		&models.Advertisement{},
		&models.Building{},
		&models.Notice{},
		&models.Administrator{},
		&models.AdvertisementBuilding{},
		&models.BuildingNotice{},
	)
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 手动添加唯一约束以防止重复关联
	err = DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_advertisement_building ON advertisement_buildings (advertisement_id, building_id);`).Error
	if err != nil {
		log.Fatal("创建唯一索引失败:", err)
	}

	// 手动添加唯一索引对于 building_notices
	err = DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_building_notice ON building_notices (building_id, notice_id);`).Error
	if err != nil {
		log.Fatal("创建 building_notices 唯一索引失败:", err)
	}
}
