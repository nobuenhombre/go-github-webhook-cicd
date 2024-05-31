package domainapp

import (
	configapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/config"
	"go-github-webhook-cicd/src/internal/pkg/services/github"
)

type AppDomain struct {
	config        *configapp.Config
	githubService github.Service
}

func NewAppDomain(config *configapp.Config) (IDomainApp, error) {
	appDomain := &AppDomain{
		config: config,
	}

	appDomain.githubService = github.NewGithub(&config.GitHub)

	return appDomain, nil
}
