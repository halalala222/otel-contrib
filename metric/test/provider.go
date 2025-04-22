package test

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

var _ metric.MeterProvider = (*Provider)(nil)
var _ metric.Meter = (*Meter)(nil)
var _ metric.Registration = (*Registration)(nil)

type Provider struct {
	embedded.MeterProvider
}

func (t *Provider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return &Meter{}
}

func (t *Provider) GetMeterData(name string, metricType string) any {
	globalMetricsMutex.Lock()
	defer globalMetricsMutex.Unlock()

	if _, ok := globalMetrics[name]; !ok {
		return nil
	}

	if _, ok := globalMetrics[name][metricType]; !ok {
		return nil
	}

	return globalMetrics[name][metricType]
}

func (t *Provider) GetMeterInt64RecordData(name string) any {
	return t.GetMeterData(name, Int64RecordType)
}

func (t *Provider) GetMeterInt64AddData(name string) any {
	return t.GetMeterData(name, Int64AddType)
}

func (t *Provider) GetMeterFloat64RecordData(name string) any {
	return t.GetMeterData(name, Float64RecordType)
}

func (t *Provider) GetMeterFloat64AddData(name string) any {
	return t.GetMeterData(name, Float64AddType)
}

type Meter struct {
	embedded.Meter
}

func (t *Meter) Int64Counter(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Int64UpDownCounter(name string, options ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Int64Histogram(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Int64Gauge(name string, options ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Int64ObservableCounter(name string, options ...metric.Int64ObservableCounterOption) (metric.Int64ObservableCounter, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Int64ObservableUpDownCounter(name string, options ...metric.Int64ObservableUpDownCounterOption) (metric.Int64ObservableUpDownCounter, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Int64ObservableGauge(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
	return &Int64MetricData{name: name}, nil
}

func (t *Meter) Float64Counter(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) Float64UpDownCounter(name string, options ...metric.Float64UpDownCounterOption) (metric.Float64UpDownCounter, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) Float64Histogram(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) Float64Gauge(name string, options ...metric.Float64GaugeOption) (metric.Float64Gauge, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) Float64ObservableCounter(name string, options ...metric.Float64ObservableCounterOption) (metric.Float64ObservableCounter, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) Float64ObservableUpDownCounter(name string, options ...metric.Float64ObservableUpDownCounterOption) (metric.Float64ObservableUpDownCounter, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) Float64ObservableGauge(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
	return &Float64MetricData{name: name}, nil
}

func (t *Meter) RegisterCallback(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
	return &Registration{}, nil
}

type Registration struct {
	embedded.Registration
}

func (t *Registration) Unregister() error {
	return nil
}
