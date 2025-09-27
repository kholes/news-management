package routes

import (
	"net/http"

	"github.com/kholes/news-management/internal/handlers"
	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, newsHandler *handlers.NewsHandler, topicHandler *handlers.TopicHandler) {
	// Root
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"message": "News Management API"})
	})

	// Topics
	e.POST("/api/topics", topicHandler.CreateTopic)
	e.GET("/api/topics", topicHandler.ListTopics)
	e.GET("/api/topics/:id", topicHandler.GetTopic)
	e.PUT("/api/topics/:id", topicHandler.UpdateTopic)
	e.DELETE("/api/topics/:id", topicHandler.DeleteTopic)

	// News
	e.POST("/api/news", newsHandler.CreateNews)
	e.GET("/api/news", newsHandler.ListNews)
	e.GET("/api/news/:id", newsHandler.GetNews)
	e.PUT("/api/news/:id", newsHandler.UpdateNews)
	e.DELETE("/api/news/:id", newsHandler.DeleteNews)
}
