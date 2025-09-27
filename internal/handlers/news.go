package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/kholes/news-management/internal/models"
)

type NewsHandler struct {
	DB *gorm.DB
}

// Constructor
func NewNewsHandler(db *gorm.DB) *NewsHandler {
	return &NewsHandler{DB: db}
}

// CreateNews godoc
// @Summary Create a new news
// @Description Create a new news with title, content, and optional topic_id
// @Tags News
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param title formData  string  true  "News Title"
// @Param content formData  string  true  "News Content"
// @Param topic_id formData  int     false "Topic ID"
// @Success 201 {object}  models.News
// @Failure 400 {object}  map[string]string
// @Failure 500 {object}  map[string]string
// @Router /news [post]
func (h *NewsHandler) CreateNews(c echo.Context) error {
	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		TopicIDs []uint `json:"topic_ids"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var topics []models.Topic
	if len(input.TopicIDs) > 0 {
		if err := h.DB.Find(&topics, input.TopicIDs).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	}

	news := models.News{
		Title:   input.Title,
		Content: input.Content,
		Status:  input.Status,
		Topics:  topics,
	}

	if err := h.DB.Create(&news).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, news)
}

// @Summary Get all news
// @Description Get list of news articles with optional filters by status and topic
// @Tags News
// @Accept json
// @Produce json
// @Param status query string false "Filter by status (draft, published, deleted)"
// @Param topic_id query int false "Filter by topic ID"
// @Success 200 {array} models.News
// @Failure 500 {object} map[string]string
// @Router /news [get]
func (h *NewsHandler) ListNews(c echo.Context) error {
	status := c.QueryParam("status")    // optional filter
	topicID := c.QueryParam("topic_id") // optional filter

	var news []models.News
	query := h.DB.Preload("Topics")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if topicID != "" {
		query = query.Joins("JOIN news_topics nt ON nt.news_id = news.id").
			Where("nt.topic_id = ?", topicID)
		query = query.Preload("Topics", "id = ?", topicID)
	} else {
		query = query.Preload("Topics")
	}

	if err := query.Find(&news).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, news)
}

// GetNews godoc
// @Summary Get news by ID
// @Description Get detailed news data including topic
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Success 200 {object} models.News
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /news/{id} [get]
func (h *NewsHandler) GetNews(c echo.Context) error {
	id := c.Param("id")
	var news models.News
	if err := h.DB.Preload("Topics").First(&news, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "News not found"})
	}
	return c.JSON(http.StatusOK, news)
}

// UpdateNews godoc
// @Summary Update an existing news
// @Description Update a news by ID with title, content, and optional topic_id (form input)
// @Tags News
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "News ID"
// @Param title formData string false "News Title"
// @Param content formData string false "News Content"
// @Param topic_id formData int false "Topic ID"
// @Success 200 {object} models.News
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /news/{id} [put]
func (h *NewsHandler) UpdateNews(c echo.Context) error {
	id := c.Param("id")
	var existing models.News
	if err := h.DB.Preload("Topics").First(&existing, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "News not found"})
	}

	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		TopicIDs []uint `json:"topic_ids"`
	}
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	existing.Title = input.Title
	existing.Content = input.Content
	existing.Status = input.Status

	if len(input.TopicIDs) > 0 {
		var topics []models.Topic
		if err := h.DB.Find(&topics, input.TopicIDs).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		h.DB.Model(&existing).Association("Topics").Replace(&topics)
	}

	if err := h.DB.Save(&existing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, existing)
}

// DeleteNews godoc
// @Summary Delete news by ID
// @Description Remove news from database
// @Tags News
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Router /news/{id} [delete]
func (h *NewsHandler) DeleteNews(c echo.Context) error {
	id := c.Param("id")
	var news models.News
	if err := h.DB.First(&news, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "News not found"})
	}

	news.Status = "deleted"
	if err := h.DB.Save(&news).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
