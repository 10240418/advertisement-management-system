package models

// Advertisement 是一个数据库模型，用于存储广告信息

import "gorm.io/gorm"

type Advertisement struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	TargetURL   string `json:"target_url"`
	Status      string `json:"status"`      // active, inactive
	BuildingID  uint   `json:"building_id"` // 外键关联 Building
	Building    Building
}
