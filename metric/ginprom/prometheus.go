package ginprom

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsPath = "/metrics"
)

// RegisterMetric exposes the Prometheus metrics endpoint.
func RegisterMetric(engin *gin.Engine) {
	engin.GET(metricsPath, func(context *gin.Context) {
		promhttp.Handler().ServeHTTP(context.Writer, context.Request)
	})
}
