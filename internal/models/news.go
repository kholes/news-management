package models

import (
	"time"

	"gorm.io/gorm"
)

type News struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	TopicID   *uint          `json:"topic_id"`
	Topic     *Topic         `json:"topic,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
