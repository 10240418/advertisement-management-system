package models

import (
	"gorm.io/gorm"
)

type Notice struct {
	gorm.Model
	Title      string `json:"title" gorm:"not null"`
	Content    string `json:"content"`
	BuildingID uint   `json:"building_id"`
	Building   Building
}
