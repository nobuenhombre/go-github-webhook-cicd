package gitexec

import (
	"github.com/nobuenhombre/suikat/pkg/ge"
	"github.com/nobuenhombre/suikat/pkg/osexec"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	"go-github-webhook-cicd/src/internal/pkg/services/queue"
	"log"
)

type Conn struct {
	config *configgithub.GitHubConfig
}

func NewGitExec(config *configgithub.GitHubConfig) Service {
	return &Conn{
		config: config,
	}
}

func (conn *Conn) exec(command string, args []string) error {
	out, err := osexec.OSRun(command, args)
	log.Printf("Exec %v \n", command)
	log.Printf("Output %v\n", out)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

func (conn *Conn) GetExecutor() queue.ExecFunc {
	return func(data interface{}) error {
		project, ok := data.(*configgithub.GitHubProjectConfig)
		if !ok {
			return ge.Pin(ge.New("data is not *configgithub.GitHubProjectConfig type"))
		}

		err := conn.exec(project.BashScripts.OnPull.Before, []string{})
		if err != nil {
			return ge.Pin(err)
		}

		err = conn.exec(conn.config.GitCmd, []string{"-C", project.Dir, "stash"})
		if err != nil {
			return ge.Pin(err)
		}

		err = conn.exec(conn.config.GitCmd, []string{"-C", project.Dir, "fetch"})
		if err != nil {
			return ge.Pin(err)
		}

		err = conn.exec(conn.config.GitCmd, []string{"-C", project.Dir, "checkout", project.Branch})
		if err != nil {
			return ge.Pin(err)
		}

		err = conn.exec(conn.config.GitCmd, []string{"-C", project.Dir, "pull"})
		if err != nil {
			return ge.Pin(err)
		}

		err = conn.exec(project.BashScripts.OnPull.After, []string{})
		if err != nil {
			return ge.Pin(err)
		}

		return nil
	}
}
