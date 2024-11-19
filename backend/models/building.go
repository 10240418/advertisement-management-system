package models

import (
	"gorm.io/gorm"
)

type Building struct {
	gorm.Model
	Name                   string                  `json:"name" gorm:"unique;not null"`
	Address                string                  `json:"address"`
	AdvertisementBuildings []AdvertisementBuilding `gorm:"foreignKey:BuildingID;constraint:OnDelete:CASCADE;" json:"advertisements_buildings"`
}
