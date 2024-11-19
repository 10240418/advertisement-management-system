package config

import (
	"fmt"
	"os"

	"github.com/10240418/advertisement-management-system/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// InitDB 初始化数据库连接，并返回错误
func InitDB() error {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return fmt.Errorf("环境变量 POSTGRES_DSN 未设置")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	err = DB.AutoMigrate(
		&models.Advertisement{},
		&models.Building{},
		&models.Administrator{},
		&models.AdvertisementBuilding{},
	)
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 手动添加唯一约束以防止重复关联
	err = DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_advertisement_building ON advertisement_buildings (advertisement_id, building_id);`).Error
	if err != nil {
		return fmt.Errorf("创建唯一索引失败: %w", err)
	}

	// 如果有其他表需要创建唯一索引，请在此处添加
	// 例如:
	// err = DB.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_building_notice ON building_notices (building_id, notice_id);`).Error
	// if err != nil {
	// 	return fmt.Errorf("创建 building_notices 唯一索引失败: %w", err)
	// }

	return nil
}
