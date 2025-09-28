package seed

import (
	"log"
	"time"

	"github.com/kholes/news-management/models"

	"gorm.io/gorm"
)

// SeedAll membuat data Topics, News, dan relasi news_topics
func SeedAll(db *gorm.DB) {
	// 1. Seed Topics
	topics := []models.Topic{
		{Name: "Technology", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "World", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Sports", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Health", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Finance", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for i := range topics {
		if err := db.FirstOrCreate(&topics[i], models.Topic{Name: topics[i].Name}).Error; err != nil {
			log.Fatalf("failed to seed topic %s: %v", topics[i].Name, err)
		}
	}

	// 2. Seed News dengan relasi otomatis ke topics (news_topics)
	newsItems := []models.News{
		{
			Title:     "Breaking News: AI takes over",
			Content:   "AI technology is transforming the world rapidly.",
			Status:    "published",
			Topics:    []models.Topic{topics[0], topics[1]}, // Technology, World
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Health Update",
			Content:   "New health research shows benefits of daily walking.",
			Status:    "draft",
			Topics:    []models.Topic{topics[3]}, // Health
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Sports Highlights",
			Content:   "Team A won the championship match yesterday.",
			Status:    "published",
			Topics:    []models.Topic{topics[2]}, // Sports
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Title:     "Finance News Today",
			Content:   "Stock markets show mixed results.",
			Status:    "published",
			Topics:    []models.Topic{topics[4]}, // Finance
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for i := range newsItems {
		if err := db.FirstOrCreate(&newsItems[i], models.News{Title: newsItems[i].Title}).Error; err != nil {
			log.Fatalf("failed to seed news %s: %v", newsItems[i].Title, err)
		}
	}
}
