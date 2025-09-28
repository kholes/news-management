package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kholes/news-management/internal/handlers"
	"github.com/kholes/news-management/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// -------------------- Setup fresh DB & handler per test --------------------
func setupTopicHandler(t *testing.T) (*handlers.TopicHandler, *gorm.DB, *echo.Echo) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	db.AutoMigrate(&models.Topic{}, &models.News{}, &models.NewsTopic{})

	h := handlers.NewTopicHandler(db)
	e := echo.New()

	return h, db, e
}

// -------------------- TEST CREATE --------------------
func TestCreateTopic(t *testing.T) {
	h, _, e := setupTopicHandler(t)

	name := fmt.Sprintf("Technology_%d", time.Now().UnixNano())
	input := map[string]string{"name": name}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/topics", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateTopic(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var topic models.Topic
	json.Unmarshal(rec.Body.Bytes(), &topic)
	assert.Equal(t, name, topic.Name)
}

// -------------------- TEST GET ALL --------------------
func TestGetAllTopics(t *testing.T) {
	h, db, e := setupTopicHandler(t)

	// Hanya insert 2 topic untuk test
	db.Create(&models.Topic{Name: "Tech"})
	db.Create(&models.Topic{Name: "Science"})

	req := httptest.NewRequest(http.MethodGet, "/api/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListTopics(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var topics []models.Topic
	json.Unmarshal(rec.Body.Bytes(), &topics)
}

// -------------------- TEST GET BY ID --------------------
func TestGetTopicByID(t *testing.T) {
	h, db, e := setupTopicHandler(t)

	topic := models.Topic{Name: "Health"}
	db.Create(&topic)

	idStr := strconv.FormatUint(uint64(topic.ID), 10)
	req := httptest.NewRequest(http.MethodGet, "/api/topics/"+idStr, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	err := h.GetTopic(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.Topic
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Health", resp.Name)
}

// -------------------- TEST UPDATE --------------------
func TestUpdateTopic(t *testing.T) {
	h, db, e := setupTopicHandler(t)

	topic := models.Topic{Name: "Old Name"}
	db.Create(&topic)

	input := map[string]string{"name": "Updated Name"}
	body, _ := json.Marshal(input)

	idStr := strconv.FormatUint(uint64(topic.ID), 10)
	req := httptest.NewRequest(http.MethodPut, "/api/topics/"+idStr, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	err := h.UpdateTopic(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.Topic
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Updated Name", resp.Name)
}

// -------------------- TEST DELETE --------------------
func TestDeleteTopic(t *testing.T) {
	h, db, e := setupTopicHandler(t)

	topic := models.Topic{Name: "Delete Me"}
	db.Create(&topic)

	idStr := strconv.FormatUint(uint64(topic.ID), 10)
	req := httptest.NewRequest(http.MethodDelete, "/api/topics/"+idStr, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	err := h.DeleteTopic(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Pastikan topik sudah dihapus
	var check models.Topic
	result := db.First(&check, topic.ID)
	assert.Error(t, result.Error)
}
