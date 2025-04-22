package gorm_plugin

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	gormDBStatusMeter = otel.Meter("gorm-db-status")
)

var (
	MaxOpenConnections, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-max-open-connections",
		metric.WithDescription("The maximum number of open connections to the database"),
	)

	// Pool status
	OpenConnections, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-open-connections",
		metric.WithDescription("The number of established connections both in use and idle"),
	)
	InUse, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-in-use",
		metric.WithDescription("The number of connections currently in use"),
	)
	Idle, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-idle",
		metric.WithDescription("The number of idle connections"),
	)

	// Counters
	WaitCount, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-wait-count",
		metric.WithDescription("The total number of connections waited for"),
	)
	WaitDuration, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-wait-duration",
		metric.WithDescription("The total time blocked waiting for a new connection"),
	)
	MaxIdleClosed, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-max-idle-closed",
		metric.WithDescription("The total number of connections closed due to SetMaxIdleConns"),
	)
	MaxLifetimeClosed, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-max-lifetime-closed",
		metric.WithDescription("The total number of connections closed due to SetConnMaxLifetime"),
	)
	MaxIdleTimeClosed, _ = gormDBStatusMeter.Float64Gauge(
		"gorm-db-status-max-idle-time-closed",
		metric.WithDescription("The total number of connections closed due to SetConnMaxIdleTime"),
	)
)
