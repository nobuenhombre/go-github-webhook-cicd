package configsentry

type SentryConfig struct {
	DSN         string `yaml:"dsn,omitempty"`
	Environment string `yaml:"environment,omitempty"`
}
