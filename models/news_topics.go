package models

type NewsTopic struct {
	NewsID  uint `gorm:"primaryKey"`
	TopicID uint `gorm:"primaryKey"`
}
