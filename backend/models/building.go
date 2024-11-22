package models

import (
	"gorm.io/gorm"
)

type Building struct {
	gorm.Model
	Name                   string                  `json:"name" gorm:"unique;not null"`
	Address                string                  `json:"address"`
	BuildingID             string                  `json:"blg_id"`
	AdvertisementBuildings []AdvertisementBuilding `gorm:"foreignKey:BuildingID;constraint:OnDelete:CASCADE;" json:"advertisements_buildings"`
}
