package handlers

import (
	"net/http"

	"github.com/kholes/news-management/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type TopicHandler struct {
	DB *gorm.DB
}

// Constructors
func NewTopicHandler(db *gorm.DB) *TopicHandler {
	return &TopicHandler{DB: db}
}

// CreateTopic godoc
// @Summary Create a new topic
// @Description  Create a new topic with name (form input)
// @Tags Topics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param name formData  string  true  "Topic Name"
// @Success 201 {object}  models.Topic
// @Failure 400 {object}  map[string]string
// @Failure 500 {object}  map[string]string
// @Router /topics [post]
func (h *TopicHandler) CreateTopic(c echo.Context) error {
	type CreateTopicInput struct {
		Name string `json:"name" binding:"required"`
	}
	var input CreateTopicInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if input.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name is required"})
	}

	topic := models.Topic{Name: input.Name}
	if err := h.DB.Create(&topic).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, topic)
}

// ListTopics godoc
// @Summary List all topics
// @Description Get list of topics
// @Tags Topics
// @Accept json
// @Produce json
// @Success 200 {array} models.Topic
// @Router /topics [get]
func (h *TopicHandler) ListTopics(c echo.Context) error {
	var topics []models.Topic
	if err := h.DB.Preload("News").Find(&topics).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, topics)
}

// GetTopic godoc
// @Summary Get topic by ID
// @Description Get detailed topic data
// @Tags Topics
// @Accept json
// @Produce json
// @Param id path int true "Topic ID"
// @Success 200 {object} models.Topic
// @Failure 404 {object} map[string]string
// @Router /topics/{id} [get]
func (h *TopicHandler) GetTopic(c echo.Context) error {
	id := c.Param("id")
	var topic models.Topic
	if err := h.DB.Preload("News").First(&topic, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "topic not found"})
	}

	return c.JSON(http.StatusOK, topic)
}

// UpdateTopic godoc
// @Summary Update an existing topic
// @Description Update a topic by ID with name (form input)
// @Tags Topics
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id path int true "Topic ID"
// @Param name formData string true "Topic Name"
// @Success 200 {object} models.Topic
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /topics/{id} [put]
func (h *TopicHandler) UpdateTopic(c echo.Context) error {
	id := c.Param("id")

	var topic models.Topic
	if err := h.DB.First(&topic, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "topic not found"})
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if input.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "name is required"})
	}

	topic.Name = input.Name
	if err := h.DB.Save(&topic).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, topic)
}

// DeleteTopic godoc
// @Summary Delete a topic
// @Description Delete a topic by its ID
// @Tags Topics
// @Param id path int  true  "Topic ID"
// @Success 204 "No Content"
// @Failure 404 {object}  map[string]string
// @Router /topics/{id} [delete]
func (h *TopicHandler) DeleteTopic(c echo.Context) error {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Topic{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
