package gorm_plugin

import "go.opentelemetry.io/otel/metric"

type prometheusConfig struct {
	dbName          string
	instrumName     string
	interval        int
	observerOptions []metric.ObserveOption
}

type Option func(*prometheusConfig)

func WithDBName(dbName string) Option {
	return func(c *prometheusConfig) {
		c.dbName = dbName
	}
}

func WithInterval(interval int) Option {
	return func(c *prometheusConfig) {
		c.interval = interval
	}
}

func WithInstrumentationName(name string) Option {
	return func(c *prometheusConfig) {
		c.instrumName = name
	}
}

func WithObserverOptions(opts ...metric.ObserveOption) Option {
	return func(c *prometheusConfig) {
		c.observerOptions = opts
	}
}

func defaultConfig() *prometheusConfig {
	return &prometheusConfig{
		dbName:   "default",
		interval: 10,
	}
}

func apply(opts ...Option) *prometheusConfig {
	newConf := defaultConfig()

	for _, opt := range opts {
		opt(newConf)
	}

	return newConf
}
