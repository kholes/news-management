package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kholes/news-management/internal/config"
	"github.com/kholes/news-management/internal/models"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (db *gorm.DB, error error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	if err := db.AutoMigrate(&models.Topic{}, &models.News{}); err != nil {
		log.Printf("AutoMigrate error: %v", err)
		return nil, err
	}

	DB = db
	return DB, nil
}

func SeedTopics(db *gorm.DB) {
	topics := []models.Topic{
		{Name: "Technology"},
		{Name: "Business"},
		{Name: "Health"},
		{Name: "Sports"},
		{Name: "Entertainment"},
	}

	for _, topic := range topics {
		var existing models.Topic
		// cek apakah sudah ada, biar tidak duplikat
		if err := db.Where("name = ?", topic.Name).First(&existing).Error; err != nil {
			if err := db.Create(&topic).Error; err != nil {
				log.Printf("Failed to seed topic %s: %v", topic.Name, err)
			} else {
				log.Printf("Seeded topic: %s", topic.Name)
			}
		}
	}
}

func SeedNews(db *gorm.DB) {
	// Pastikan sudah ada topics dulu
	var topics []models.Topic
	if err := DB.Find(&topics).Error; err != nil {
		log.Printf("Failed to fetch topics: %v", err)
		return
	}
	if len(topics) == 0 {
		log.Println("No topics found. Run SeedTopics() first")
		return
	}

	newsList := []models.News{
		{Title: "Go 1.24 Released", Content: "The Go team just released Go 1.24 with exciting new features.", TopicID: &topics[0].ID},
		{Title: "Stock Market Update", Content: "The stock market saw significant growth today.", TopicID: &topics[1].ID},
		{Title: "Health Tips 2025", Content: "Experts suggest daily exercise to improve mental health.", TopicID: &topics[2].ID},
		{Title: "World Cup 2026 Preparations", Content: "Countries are preparing stadiums and teams for World Cup 2026.", TopicID: &topics[3].ID},
		{Title: "New Marvel Movie Released", Content: "Marvel Studios released its latest superhero blockbuster.", TopicID: &topics[4].ID},
	}

	for _, n := range newsList {
		var existing models.News
		if err := db.Where("title = ?", n.Title).First(&existing).Error; err != nil {
			if err := db.Create(&n).Error; err != nil {
				log.Printf("Failed to seed news %s: %v", n.Title, err)
			} else {
				log.Printf("Seeded news: %s", n.Title)
			}
		}
	}
}
