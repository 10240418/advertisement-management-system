package models

// Advertisement 是一个数据库模型，用于存储广告信息

import (
	"gorm.io/gorm"
)

type Advertisement struct {
	gorm.Model
	Title                  string                  `json:"title" binding:"required"`
	Description            string                  `json:"description"`
	ImageURL               string                  `json:"image_url"`
	VideoURL               string                  `json:"video_url"`
	VideoDuration          int64                   `json:"video_duration"`            // 以秒为单位
	Status                 string                  `json:"status" binding:"required"` // active, inactive
	AdvertisementBuildings []AdvertisementBuilding `gorm:"foreignKey:AdvertisementID" json:"advertisements_buildings"`
}
