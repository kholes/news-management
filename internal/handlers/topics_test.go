package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/kholes/news-management/internal/handlers"
	"github.com/kholes/news-management/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func topicsSetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to sqlite in memory: %v", err)
	}

	// migrate ulang
	db.AutoMigrate(&models.News{}, &models.Topic{})

	// bersihkan tabel
	db.Exec("DELETE FROM news")
	db.Exec("DELETE FROM topics")

	return db
}

func TestCreateTopics(t *testing.T) {
	e := echo.New()
	db := topicsSetupTestDB(t)
	h := handlers.NewTopicHandler(db)

	// form data
	form := make(url.Values)
	form.Set("name", "Tech")

	req := httptest.NewRequest(http.MethodPost, "/topics", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm) // application/x-www-form-urlencoded
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call handler
	if assert.NoError(t, h.CreateTopic(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.Topic
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "Tech", resp.Name)
		assert.NotZero(t, resp.ID)
	}
}

func TestGetAll(t *testing.T) {
	e := echo.New()
	db := topicsSetupTestDB(t)

	// pastikan DB kosong setelah test selesai
	t.Cleanup(func() {
		db.Exec("DELETE FROM news")
		db.Exec("DELETE FROM topics")
	})

	h := handlers.NewTopicHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/topics", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.ListTopics(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetTopic_Success(t *testing.T) {
	e := echo.New()
	db := topicsSetupTestDB(t)
	h := handlers.NewTopicHandler(db)

	// insert dummy topic
	topic := models.Topic{Name: "Technology"}
	db.Create(&topic)

	req := httptest.NewRequest(http.MethodGet, "/topics/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/topics/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.GetTopic(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Technology")
	}
}

func TestGetTopic_NotFound(t *testing.T) {
	e := echo.New()
	db := topicsSetupTestDB(t)
	h := handlers.NewTopicHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/topics/99", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/topics/:id")
	c.SetParamNames("id")
	c.SetParamValues("99")

	if assert.NoError(t, h.GetTopic(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "topic not found")
	}
}

func TestUpdateTopics_Success(t *testing.T) {
	e := echo.New()
	db := topicsSetupTestDB(t)
	h := handlers.NewTopicHandler(db)

	// seed
	topic := models.Topic{Name: "Old Topic"}
	db.Create(&topic)

	form := url.Values{}
	form.Set("name", "Updated Topic")

	req := httptest.NewRequest(http.MethodPut, "/topics/"+strconv.Itoa(int(topic.ID)), strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(topic.ID)))

	if assert.NoError(t, h.UpdateTopic(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updated models.Topic
		json.Unmarshal(rec.Body.Bytes(), &updated)
		assert.Equal(t, "Updated Topic", updated.Name)
	}
}

func TestUpdateTopics_NotFound(t *testing.T) {
	e := echo.New()
	db := topicsSetupTestDB(t)
	h := handlers.NewTopicHandler(db)

	// tidak seed data → news id=99 tidak ada
	payload := models.News{Title: "Does Not Matter"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/topics/99", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/topics/:id")
	c.SetParamNames("id")
	c.SetParamValues("99")

	// call handler
	if assert.NoError(t, h.UpdateTopic(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}
