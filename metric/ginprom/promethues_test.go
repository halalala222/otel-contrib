package ginprom

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestRegisterMetric(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	RegisterMetric(engine)

	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()

	engine.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "text/plain; version=0.0.4; charset=utf-8; escaping=underscores", response.Header().Get("Content-Type"))
	assert.NotEqual(t, 0, len(response.Body.String()))
}
