package domainapp

import "go-github-webhook-cicd/src/internal/pkg/services/github"

type IDomainApp interface {
	GetGithubService() github.Service
}
