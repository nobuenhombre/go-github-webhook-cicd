package domainapp

import (
	"go-github-webhook-cicd/src/internal/pkg/services/github"
	"go-github-webhook-cicd/src/internal/pkg/services/queue"
)

type IDomainApp interface {
	GetGithubService() github.Service
	GetQueueService() queue.Service
}
