package github

import (
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	"net/http"
)

type Service interface {
	GetProjects() []configgithub.GitHubProjectConfig
	ValidatePushEventRequest(request *http.Request, project *configgithub.GitHubProjectConfig) error
	OnPush() error
}
