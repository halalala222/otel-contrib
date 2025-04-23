package gormprom

import "go.opentelemetry.io/otel/metric"

const defaultInstrumentName = "ncuhome:gorm"

type prometheusConfig struct {
	instrumName     string
	observerOptions []metric.ObserveOption
}

type Option func(*prometheusConfig)

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
		instrumName: defaultInstrumentName,
	}
}

func apply(opts ...Option) *prometheusConfig {
	newConf := defaultConfig()

	for _, opt := range opts {
		opt(newConf)
	}

	return newConf
}
