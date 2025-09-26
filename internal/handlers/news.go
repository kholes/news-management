package handlers

import (
	"net/http"
	"strconv"

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
	title := c.FormValue("title")
	content := c.FormValue("content")

	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is required"})
	}
	if content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "content is required"})
	}

	var topicID *uint
	if topicIDStr := c.FormValue("topic_id"); topicIDStr != "" {
		idUint, err := strconv.ParseUint(topicIDStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid topic_id"})
		}
		id := uint(idUint)
		topicID = &id
	}

	news := models.News{Title: title, Content: content, TopicID: topicID}
	if err := h.DB.Create(&news).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, news)
}

// ListNews godoc
// @Summary List all news
// @Description Get list of news
// @Tags News
// @Accept json
// @Produce json
// @Success 200 {array} models.News
// @Router /news [get]
func (h *NewsHandler) ListNews(c echo.Context) error {
	var news []models.News
	limit := 20
	offset := 0
	if q := c.QueryParam("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = v
		}
	}
	if q := c.QueryParam("offset"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			offset = v
		}
	}

	// optionally filter by topic
	if q := c.QueryParam("topic_id"); q != "" {
		if err := h.DB.Preload("Topic").Where("topic_id = ?", q).Limit(limit).Offset(offset).Find(&news).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, news)
	}

	if err := h.DB.Preload("Topic").Limit(limit).Offset(offset).Find(&news).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
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
	var n models.News
	if err := h.DB.Preload("Topic").First(&n, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "news not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, n)
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

	var news models.News
	if err := h.DB.First(&news, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "news not found"})
	}

	if title := c.FormValue("title"); title != "" {
		news.Title = title
	}
	if content := c.FormValue("content"); content != "" {
		news.Content = content
	}
	if topicIDStr := c.FormValue("topic_id"); topicIDStr != "" {
		idUint, err := strconv.ParseUint(topicIDStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid topic_id"})
		}
		id := uint(idUint)
		news.TopicID = &id
	}

	if err := h.DB.Save(&news).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, news)
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
	if err := h.DB.Delete(&models.News{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
