package gorm_plugin

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var _ gorm.Plugin = &Prometheus{}

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

	p.refreshOnce.Do(func() {
		go func() {
			for range time.Tick(time.Duration(p.config.interval) * time.Second) {
				p.refresh()
			}
		}()
	})

	p.refresh()

	return nil
}

func New(options ...Option) *Prometheus {
	return &Prometheus{
		config: apply(options...),
	}
}

func (p *Prometheus) refreshStats(dbStats sql.DBStats) {
	MaxOpenConnections.Record(context.Background(), float64(dbStats.MaxOpenConnections))
	OpenConnections.Record(context.Background(), float64(dbStats.OpenConnections))
	InUse.Record(context.Background(), float64(dbStats.InUse))
	Idle.Record(context.Background(), float64(dbStats.Idle))
	WaitCount.Record(context.Background(), float64(dbStats.WaitCount))
	WaitDuration.Record(context.Background(), float64(dbStats.WaitDuration))
	MaxIdleClosed.Record(context.Background(), float64(dbStats.MaxIdleClosed))
	MaxLifetimeClosed.Record(context.Background(), float64(dbStats.MaxLifetimeClosed))
	MaxIdleTimeClosed.Record(context.Background(), float64(dbStats.MaxIdleTimeClosed))
}

func (p *Prometheus) refresh() {
	var (
		db  *sql.DB
		err error
	)

	if db, err = p.DB.DB(); err != nil {
		p.DB.Logger.Error(context.Background(), "gorm:ncuhome:otel:prometheus failed to refresh db stats,got error: %v", err)
		return
	}

	p.refreshStats(db.Stats())
}
