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

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kholes/news-management/internal/handlers"
	"github.com/kholes/news-management/internal/models"
)

func newsSetupTestDB(t *testing.T) *gorm.DB {
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

func TestCreateNews_Form(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)

	form := strings.NewReader("title=Breaking News&content=AI takes over")
	req := httptest.NewRequest(http.MethodPost, "/news", form)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm) // penting
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.CreateNews(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp models.News
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, "Breaking News", resp.Title)
		assert.Equal(t, "AI takes over", resp.Content)
		assert.NotZero(t, resp.ID)
	}
}

func TestGetAllNews(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)

	// pastikan DB kosong setelah test selesai
	t.Cleanup(func() {
		db.Exec("DELETE FROM news")
		db.Exec("DELETE FROM topics")
	})

	h := handlers.NewNewsHandler(db)

	// Insert data manual ke DB
	db.Create(&models.News{Title: "News A", Content: "Content A"})
	db.Create(&models.News{Title: "News B", Content: "Content B"})

	req := httptest.NewRequest(http.MethodGet, "/news", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, h.ListNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []models.News
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Len(t, resp, 2)
		assert.Equal(t, "News A", resp[0].Title)
	}
}

func TestUpdateNews_Success(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)

	// seed
	news := models.News{Title: "Old Title", Content: "Old Content"}
	db.Create(&news)

	form := url.Values{}
	form.Set("title", "Updated Title")
	form.Set("content", "Updated Content")

	req := httptest.NewRequest(http.MethodPut, "/news/"+strconv.Itoa(int(news.ID)), strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(news.ID)))

	if assert.NoError(t, h.UpdateNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var updated models.News
		json.Unmarshal(rec.Body.Bytes(), &updated)
		assert.Equal(t, "Updated Title", updated.Title)
	}
}

func TestUpdateNews_NotFound(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)

	// tidak seed data → news id=99 tidak ada
	payload := models.News{Title: "Does Not Matter"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/news/99", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/news/:id")
	c.SetParamNames("id")
	c.SetParamValues("99")

	// call handler
	if assert.NoError(t, h.UpdateNews(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestDeleteNews_Success(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)

	// seed news
	news := models.News{Title: "To be deleted", Content: "Bye"}
	db.Create(&news)

	req := httptest.NewRequest(http.MethodDelete, "/news/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/news/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	// call handler
	if assert.NoError(t, h.DeleteNews(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// pastikan news sudah hilang
		var count int64
		db.Model(&models.News{}).Where("id = ?", 1).Count(&count)
		assert.Equal(t, int64(0), count)
	}
}

func TestDeleteNews_NotFound(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)
	req := httptest.NewRequest(http.MethodDelete, "/news/99", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/news/:id")
	c.SetParamNames("id")
	c.SetParamValues("99")

	// call handler
	if assert.NoError(t, h.DeleteNews(c)) {
		// tetap return 204 meskipun record tidak ada (sesuai default GORM)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestGetNews_Success(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)

	// insert dummy topic + news
	topic := models.Topic{Name: "Tech"}
	db.Create(&topic)

	news := models.News{
		Title:   "Test News",
		Content: "This is content",
		TopicID: &topic.ID,
	}
	db.Create(&news)

	req := httptest.NewRequest(http.MethodGet, "/news/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/news/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, h.GetNews(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Test News")
		assert.Contains(t, rec.Body.String(), "Tech")
	}
}

func TestGetNews_NotFound(t *testing.T) {
	e := echo.New()
	db := newsSetupTestDB(t)
	h := handlers.NewNewsHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/news/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/news/:id")
	c.SetParamNames("id")
	c.SetParamValues("999")

	if assert.NoError(t, h.GetNews(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "news not found")
	}
}
