package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/ncuhome/otel-contrib/metric/ginprom"
	"github.com/ncuhome/otel-contrib/metric/gormprom"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGinWithGorm(t *testing.T) {
	ctx := context.Background()
	engine := gin.New()

	ginprom.RegisterMetric(engine)

	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))

	if err != nil {
		t.Fatal(err)
	}

	if err = db.Use(gormprom.New()); err != nil {
		t.Fatal(err)
	}

	var num int64
	if err = db.WithContext(ctx).Raw("SELECT 32").Scan(&num).Error; err != nil {
		t.Fatal(err)
	}

	go func() {
		if err = engine.Run(":8080"); err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(1 * time.Second)

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	engine.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "text/plain; version=0.0.4; charset=utf-8; escaping=underscores", response.Header().Get("Content-Type"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_max_open"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_open"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_in_use"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_idle"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_wait_count"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_wait_duration"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_closed_max_idle"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_closed_max_idle_time"))
	assert.Equal(t, true, strings.Contains(response.Body.String(), "go_sql_connections_closed_max_lifetime"))
	assert.NotEqual(t, 0, len(response.Body.String()))
}
