package models

import (
	"gorm.io/gorm"
)

type Building struct {
	gorm.Model
	Name           string          `json:"name" gorm:"unique;not null"`
	Address        string          `json:"address"`
	Advertisements []Advertisement `json:"advertisements"`
	Notices        []Notice        `json:"notices"`
}
