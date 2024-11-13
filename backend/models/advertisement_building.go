package models

// AdvertisementBuilding 是广告与大厦之间的关联模型，用于存储播放时长
type AdvertisementBuilding struct {
	AdvertisementID uint          `json:"advertisement_id" gorm:"not null"`
	BuildingID      uint          `json:"building_id" gorm:"not null"`
	PlayDuration    int64         `json:"play_duration"` // 以秒为单位
	Advertisement   Advertisement `gorm:"foreignKey:AdvertisementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Building        Building      `gorm:"foreignKey:BuildingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName 设置表名
func (AdvertisementBuilding) TableName() string {
	return "advertisement_buildings"
}
