package models

import (
	"gorm.io/gorm"
)

type Administrator struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"type:varchar(100);not null"`
}
