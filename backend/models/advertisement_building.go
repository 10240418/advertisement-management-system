package models

// AdvertisementBuilding 是广告与大厦之间的关联模型，用于存储播放时长
type AdvertisementBuilding struct {
	AdvertisementID uint          `json:"advertisement_id"`
	BuildingID      uint          `json:"building_id"`
	PlayDuration    int64         `json:"play_duration"` // 以秒为单位
	Advertisement   Advertisement `gorm:"foreignKey:AdvertisementID"`
	Building        Building      `gorm:"foreignKey:BuildingID"`
}
