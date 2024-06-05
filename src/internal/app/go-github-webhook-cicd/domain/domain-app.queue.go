package domainapp

import (
	"go-github-webhook-cicd/src/internal/pkg/services/queue"
)

func (c *AppDomain) GetQueueService() queue.Service {
	return c.queueService
}
