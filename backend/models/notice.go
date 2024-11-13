package models

import (
	"gorm.io/gorm"
)

type Notice struct {
	gorm.Model
	Title     string     `json:"title" gorm:"not null"`
	PDFURL    string     `json:"pdf_url"`
	Buildings []Building `gorm:"many2many:building_notices;" json:"building_ids"`
}

type BuildingNotice struct {
	BuildingID uint
	NoticeID   uint
}
