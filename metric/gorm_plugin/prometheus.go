package gorm_plugin

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"gorm.io/gorm"
)

var _ gorm.Plugin = &Prometheus{}

const (
	maxOpenConnsName           = "go.sql.connections_max_open"
	openConnsName              = "go.sql.connections_open"
	inUseConnsName             = "go.sql.connections_in_use"
	idleConnsName              = "go.sql.connections_idle"
	connsWaitCountName         = "go.sql.connections_wait_count"
	connsWaitDurationName      = "go.sql.connections_wait_duration"
	connsClosedMaxIdleName     = "go.sql.connections_closed_max_idle"
	connsClosedMaxIdleTimeName = "go.sql.connections_closed_max_idle_time"
	connsClosedMaxLifetimeName = "go.sql.connections_closed_max_lifetime"
)

type Prometheus struct {
	*gorm.DB
	config      *prometheusConfig
	refreshOnce sync.Once
}

func (p *Prometheus) Name() string {
	return "gorm:ncuhome:otel:prometheus"
}

func (p *Prometheus) Initialize(db *gorm.DB) error {
	p.DB = db
	meter := otel.Meter(p.config.instrumName)

	maxOpenConns, _ := meter.Int64ObservableGauge(
		maxOpenConnsName,
		metric.WithDescription("Maximum number of open connections to the database"),
	)
	openConns, _ := meter.Int64ObservableGauge(
		openConnsName,
		metric.WithDescription("The number of established connections both in use and idle"),
	)
	inUseConns, _ := meter.Int64ObservableGauge(
		inUseConnsName,
		metric.WithDescription("The number of connections currently in use"),
	)
	idleConns, _ := meter.Int64ObservableGauge(
		idleConnsName,
		metric.WithDescription("The number of idle connections"),
	)
	connsWaitCount, _ := meter.Int64ObservableCounter(
		connsWaitCountName,
		metric.WithDescription("The total number of connections waited for"),
	)
	connsWaitDuration, _ := meter.Int64ObservableCounter(
		connsWaitDurationName,
		metric.WithDescription("The total time blocked waiting for a new connection"),
		metric.WithUnit("nanoseconds"),
	)
	connsClosedMaxIdle, _ := meter.Int64ObservableCounter(
		connsClosedMaxIdleName,
		metric.WithDescription("The total number of connections closed due to SetMaxIdleConns"),
	)
	connsClosedMaxIdleTime, _ := meter.Int64ObservableCounter(
		connsClosedMaxIdleTimeName,
		metric.WithDescription("The total number of connections closed due to SetConnMaxIdleTime"),
	)
	connsClosedMaxLifetime, _ := meter.Int64ObservableCounter(
		connsClosedMaxLifetimeName,
		metric.WithDescription("The total number of connections closed due to SetConnMaxLifetime"),
	)

	_, err := meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			sqlDB, err := p.DB.DB()

			if err != nil {
				return err
			}

			stats := sqlDB.Stats()

			o.ObserveInt64(maxOpenConns, int64(stats.MaxOpenConnections), p.config.observerOptions...)
			o.ObserveInt64(openConns, int64(stats.OpenConnections), p.config.observerOptions...)
			o.ObserveInt64(inUseConns, int64(stats.InUse), p.config.observerOptions...)
			o.ObserveInt64(idleConns, int64(stats.Idle), p.config.observerOptions...)
			o.ObserveInt64(connsWaitCount, stats.WaitCount, p.config.observerOptions...)
			o.ObserveInt64(connsWaitDuration, int64(stats.WaitDuration), p.config.observerOptions...)
			o.ObserveInt64(connsClosedMaxIdle, stats.MaxIdleClosed, p.config.observerOptions...)
			o.ObserveInt64(connsClosedMaxIdleTime, stats.MaxIdleTimeClosed, p.config.observerOptions...)
			o.ObserveInt64(connsClosedMaxLifetime, stats.MaxLifetimeClosed, p.config.observerOptions...)
			return nil
		},
		maxOpenConns,
		openConns,
		inUseConns,
		idleConns,
		connsWaitCount,
		connsWaitDuration,
		connsClosedMaxIdle,
		connsClosedMaxIdleTime,
		connsClosedMaxLifetime,
	)

	return err
}

func New(options ...Option) *Prometheus {
	return &Prometheus{
		config: apply(options...),
	}
}
