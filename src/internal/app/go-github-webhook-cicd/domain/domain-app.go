package domainapp

import (
	configapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/config"
	gitexec "go-github-webhook-cicd/src/internal/pkg/services/gitExec"
	"go-github-webhook-cicd/src/internal/pkg/services/github"
	"go-github-webhook-cicd/src/internal/pkg/services/queue"
	"time"
)

type AppDomain struct {
	config         *configapp.Config
	githubService  github.Service
	gitExecService gitexec.Service
	queueService   queue.Service
}

func NewAppDomain(config *configapp.Config) (IDomainApp, error) {
	appDomain := &AppDomain{
		config: config,
	}

	appDomain.githubService = github.NewGithub(&config.GitHub)
	appDomain.gitExecService = gitexec.NewGitExec(&config.GitHub)
	appDomain.queueService = queue.NewQueue(appDomain.gitExecService.GetExecutor(), 500*time.Millisecond)

	return appDomain, nil
}
