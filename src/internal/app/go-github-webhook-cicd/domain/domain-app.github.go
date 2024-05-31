package domainapp

import "go-github-webhook-cicd/src/internal/pkg/services/github"

func (c *AppDomain) GetGithubService() github.Service {
	return c.githubService
}
