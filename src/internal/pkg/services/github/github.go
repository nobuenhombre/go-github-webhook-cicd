package github

import configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"

type Conn struct {
	config *configgithub.GitHubConfig
}

func NewGithub(config *configgithub.GitHubConfig) Service {
	return &Conn{
		config: config,
	}
}

func (c *Conn) Pull() error {
	return nil
}
