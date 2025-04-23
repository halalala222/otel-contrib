package test

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

const (
	Int64RecordType = "int64_record"
	Int64AddType    = "int64_add"

	Float64RecordType = "float64_record"
	Float64AddType    = "float64_add"
)

var (
	globalMetrics      = make(map[string]map[string]any)
	globalMetricsMutex = &sync.Mutex{}
	writeMetricValue   = &writeMetric{}
)

var _ metric.Int64Counter = (*Int64MetricData)(nil)
var _ metric.Int64UpDownCounter = (*Int64MetricData)(nil)
var _ metric.Int64Histogram = (*Int64MetricData)(nil)
var _ metric.Int64Gauge = (*Int64MetricData)(nil)
var _ metric.Int64ObservableCounter = (*Int64MetricData)(nil)
var _ metric.Int64ObservableUpDownCounter = (*Int64MetricData)(nil)
var _ metric.Int64ObservableGauge = (*Int64MetricData)(nil)

var _ metric.Float64Counter = (*Float64MetricData)(nil)
var _ metric.Float64UpDownCounter = (*Float64MetricData)(nil)
var _ metric.Float64Histogram = (*Float64MetricData)(nil)
var _ metric.Float64Gauge = (*Float64MetricData)(nil)
var _ metric.Float64ObservableCounter = (*Float64MetricData)(nil)
var _ metric.Float64ObservableUpDownCounter = (*Float64MetricData)(nil)
var _ metric.Float64ObservableGauge = (*Float64MetricData)(nil)

var _ metric.Observer = (*Observer)(nil)

// Int64MetricData is a test implementation of the Int64 type of metric.
type Int64MetricData struct {
	metric.Int64Observable
	embedded.Int64Counter
	embedded.Int64UpDownCounter
	embedded.Int64Histogram
	embedded.Int64Gauge
	embedded.Int64Observer
	embedded.Int64ObservableCounter
	embedded.Int64ObservableUpDownCounter
	embedded.Int64ObservableGauge

	name string
}

type writeMetric struct {
}

func (t *writeMetric) Write(metricName string, key string, value any) {
	globalMetricsMutex.Lock()
	defer globalMetricsMutex.Unlock()

	if _, ok := globalMetrics[metricName]; !ok {
		globalMetrics[metricName] = make(map[string]any)
	}

	globalMetrics[metricName][key] = value
}

func (t *Int64MetricData) Record(ctx context.Context, incr int64, options ...metric.RecordOption) {
	writeMetricValue.Write(t.name, Int64RecordType, incr)
}

func (t *Int64MetricData) Add(ctx context.Context, incr int64, options ...metric.AddOption) {
	writeMetricValue.Write(t.name, Int64AddType, incr)
}

// Float64MetricData is a test implementation of the Float64 type of metric.
type Float64MetricData struct {
	metric.Float64Observable
	embedded.Float64Counter
	embedded.Float64UpDownCounter
	embedded.Float64Histogram
	embedded.Float64Gauge
	embedded.Float64Observer
	embedded.Float64ObservableCounter
	embedded.Float64ObservableUpDownCounter
	embedded.Float64ObservableGauge

	name string
}

func (t *Float64MetricData) Record(ctx context.Context, incr float64, options ...metric.RecordOption) {
	writeMetricValue.Write(t.name, Float64RecordType, incr)
}

func (t *Float64MetricData) Add(ctx context.Context, incr float64, options ...metric.AddOption) {
	writeMetricValue.Write(t.name, Float64AddType, incr)
}

type Observer struct {
	embedded.Observer
}

func (o *Observer) ObserveInt64(obsrv metric.Int64Observable, value int64, opts ...metric.ObserveOption) {
	if obsrv == nil {
		return
	}

	if _, ok := obsrv.(*Int64MetricData); !ok {
		return
	}

	int64MetricData := obsrv.(*Int64MetricData)
	int64MetricData.Record(context.Background(), value)
}

func (o *Observer) ObserveFloat64(obsrv metric.Float64Observable, value float64, opts ...metric.ObserveOption) {
	if obsrv == nil {
		return
	}

	if _, ok := obsrv.(Float64MetricData); !ok {
		return
	}

	float64MetricData := obsrv.(Float64MetricData)
	float64MetricData.Record(context.Background(), value)
}
