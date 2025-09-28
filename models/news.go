package models

import (
	"time"
)

type News struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Status    string    `gorm:"type:varchar(20);not null;check:status IN ('draft','published','deleted')" json:"status"`
	Topics    []Topic   `gorm:"many2many:news_topics;" json:"topics,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
