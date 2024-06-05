package gitexec

import (
	"go-github-webhook-cicd/src/internal/pkg/services/queue"
)

type Service interface {
	GetExecutor() queue.ExecFunc
}
