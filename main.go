package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/kholes/news-management/docs" // generated docs

	"github.com/kholes/news-management/config"
	"github.com/kholes/news-management/database"
	"github.com/kholes/news-management/internal/handlers"
	"github.com/kholes/news-management/internal/routes"
	"github.com/kholes/news-management/models"
	seed "github.com/kholes/news-management/seeders"
)

// @title News and Topics Management API
// @version 1.0
// @description API for managing news and topics
// @host localhost:8080
// @BasePath /api

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed connect database: %v", err)
	}

	// Migrate
	db.AutoMigrate(&models.News{}, &models.Topic{}, &models.NewsTopic{})

	// Seed
	seed.SeedAll(db)

	newsHandler := handlers.NewNewsHandler(db)
	topicHandler := handlers.NewTopicHandler(db)

	e := echo.New()
	routes.Register(e, newsHandler, topicHandler)

	// Swagger docs
	e.GET("/docs/*", echoSwagger.WrapHandler)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on %s", addr)
	e.Logger.Fatal(e.Start(addr))
}
