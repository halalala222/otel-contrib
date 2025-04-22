package gorm_plugin

type prometheusConfig struct {
	dbName        string
	interval      int
	variableNames []string
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
