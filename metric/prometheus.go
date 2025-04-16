package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	metricsPath = "/metrics"
)

func init() {
	exporter, err := prometheus.New()

	if err != nil {
		panic(err)
	}

	otel.SetMeterProvider(sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter)))
}

// RegisterMetric exposes the Prometheus metrics endpoint.
func RegisterMetric(engin *gin.Engine) {
	engin.GET(metricsPath, func(context *gin.Context) {
		promhttp.Handler().ServeHTTP(context.Writer, context.Request)
	})
}
