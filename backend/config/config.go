package config

import (
	"log"
	"os"

	"github.com/10240418/advertisement-management-system/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("POSTGRES_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移模型
	err = DB.AutoMigrate(&models.Advertisement{}, &models.Building{}, &models.Notice{}, &models.Administrator{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
