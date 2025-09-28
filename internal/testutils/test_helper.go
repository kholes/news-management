package testutils

import (
	"log"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kholes/news-management/models"
)

// SetupTestDB buat DB in-memory
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect test db: %v", err)
	}

	err = db.AutoMigrate(&models.News{}, &models.Topic{}, &models.NewsTopic{})
	if err != nil {
		t.Fatalf("failed to migrate db: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM news_topics")
		db.Exec("DELETE FROM news")
		db.Exec("DELETE FROM topics")
	})

	return db
}

// Seed topics
func SeedTopics(db *gorm.DB, names []string) []models.Topic {
	topics := []models.Topic{}
	for _, n := range names {
		t := models.Topic{Name: n, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		db.Create(&t)
		topics = append(topics, t)
	}
	return topics
}

// Seed news
func SeedNews(db *gorm.DB, newsList []models.News) []models.News {
	for i := range newsList {
		db.Create(&newsList[i])
	}
	return newsList
}
