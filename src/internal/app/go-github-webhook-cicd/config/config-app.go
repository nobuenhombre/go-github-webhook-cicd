package configapp

import (
	"github.com/nobuenhombre/suikat/pkg/fico"
	"github.com/nobuenhombre/suikat/pkg/ge"
	configserver "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/config"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	configsentry "go-github-webhook-cicd/src/internal/pkg/services/sentry/config"
	"gopkg.in/yaml.v3"
)

type HostsConfig struct {
	API configserver.HTTPServerConfig `yaml:"api,omitempty"`
}

type Config struct {
	Hosts  HostsConfig               `yaml:"hosts,omitempty"`
	GitHub configgithub.GitHubConfig `yaml:"github,omitempty"`
	Sentry configsentry.SentryConfig `yaml:"sentry,omitempty"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(fileName string) error {
	txtConfigFile := fico.TxtFile(fileName)

	configData, err := txtConfigFile.Read()
	if err != nil {
		return ge.Pin(err)
	}

	err = yaml.Unmarshal([]byte(configData), c)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

func (c *Config) Save(fileName string) error {
	txtConfigFile := fico.TxtFile(fileName)

	configData, err := yaml.Marshal(c)
	if err != nil {
		return ge.Pin(err)
	}

	err = txtConfigFile.Write(string(configData))
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
