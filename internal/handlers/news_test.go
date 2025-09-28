package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/kholes/news-management/internal/handlers"
	"github.com/kholes/news-management/models"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupNewsHandler(t *testing.T) (*handlers.NewsHandler, *gorm.DB, *echo.Echo) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	db.AutoMigrate(&models.News{}, &models.Topic{}, &models.NewsTopic{})

	h := handlers.NewNewsHandler(db)
	e := echo.New()

	return h, db, e
}

// -------------------- TEST CREATE --------------------
func TestCreateNews(t *testing.T) {
	h, db, e := setupNewsHandler(t)

	// Seed topic
	topic := models.Topic{Name: "Technology"}
	db.Create(&topic)

	input := map[string]interface{}{
		"title":     "Breaking News",
		"content":   "AI takes over",
		"status":    "draft",
		"topic_ids": []uint{topic.ID},
	}
	body, _ := json.Marshal(input)

	req := httptest.NewRequest(http.MethodPost, "/api/news", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateNews(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

// -------------------- TEST GET BY ID --------------------
func TestGetNewsByID(t *testing.T) {
	h, db, e := setupNewsHandler(t)

	// Seed News
	news := models.News{
		Title:     "Test News",
		Content:   "Content",
		Status:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Create(&news)

	req := httptest.NewRequest(http.MethodGet, "/api/news/"+strconv.FormatUint(uint64(news.ID), 10), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatUint(uint64(news.ID), 10))

	err := h.GetNews(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.News
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, news.Title, resp.Title)
}

// -------------------- TEST UPDATE --------------------
func TestUpdateNews(t *testing.T) {
	h, db, e := setupNewsHandler(t)

	// Seed News
	news := models.News{
		Title:   "Old Title",
		Content: "Old Content",
		Status:  "draft",
	}
	db.Create(&news)

	input := map[string]interface{}{
		"title":     "Updated Title",
		"content":   "Updated Content",
		"status":    "published",
		"topic_ids": []uint{},
	}
	body, _ := json.Marshal(input)

	idStr := strconv.FormatUint(uint64(news.ID), 10)
	req := httptest.NewRequest(http.MethodPut, "/api/news/"+idStr, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	err := h.UpdateNews(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.News
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "Updated Title", resp.Title)
	assert.Equal(t, "Updated Content", resp.Content)
	assert.Equal(t, "published", resp.Status)
}

// -------------------- TEST DELETE --------------------
func TestDeleteNews(t *testing.T) {
	h, db, e := setupNewsHandler(t)

	// Seed News
	news := models.News{
		Title:   "Delete Me",
		Content: "Some content",
		Status:  "draft",
	}
	db.Create(&news)

	idStr := strconv.FormatUint(uint64(news.ID), 10)
	req := httptest.NewRequest(http.MethodDelete, "/api/news/"+idStr, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(idStr)

	err := h.DeleteNews(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Pastikan news status diubah menjadi "deleted"
	var check models.News
	result := db.First(&check, news.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, "deleted", check.Status)
}
