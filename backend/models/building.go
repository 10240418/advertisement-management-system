package models

import (
	"gorm.io/gorm"
)

type Building struct {
	gorm.Model
	Name                   string                  `json:"name" gorm:"unique;not null"`
	Address                string                  `json:"address"`
	AdvertisementBuildings []AdvertisementBuilding `gorm:"foreignKey:BuildingID" json:"advertisements_buildings"`
	BuildingNotices        []BuildingNotice        `gorm:"foreignKey:BuildingID" json:"building_notices"`
}
