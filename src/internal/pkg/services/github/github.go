package github

import (
	"github.com/nobuenhombre/suikat/pkg/ge"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	"net/http"
)

type Conn struct {
	config *configgithub.GitHubConfig
}

func NewGithub(config *configgithub.GitHubConfig) Service {
	return &Conn{
		config: config,
	}
}

func (c *Conn) GetProjects() []configgithub.GitHubProjectConfig {
	return c.config.Projects
}

func (c *Conn) OnPush() error {
	return nil
}

func (c *Conn) ValidatePushEventRequest(r *http.Request, project *configgithub.GitHubProjectConfig) error {
	peRequest, err := NewPushEventRequest(r)
	if err != nil {
		return ge.Pin(err)
	}

	err = peRequest.Validate(project.Secret, project.Branch)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
