package cli

import (
	"github.com/nobuenhombre/suikat/pkg/clivar"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

// Config App config
type Config struct {
	ConfigFile string `cli:"config[config file]:string=config.api.yaml"`
	LogFile    string `cli:"log[log file]:string=/var/log/go-github-webhook-cicd/api.log"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := clivar.Load(cfg)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return cfg, nil
}
