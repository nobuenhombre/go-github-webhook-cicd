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

func (conn *Conn) exec(ident string, command string, args []string) error {
	out, err := osexec.OSRun(command, args)
	log.Printf("[%v] Exec %v args %v \n", ident, command, args)
	log.Printf("[%v] \tOutput: %v\n", ident, out)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

type GitCommand struct {
	Cmd  string
	Args []string
}

func (conn *Conn) GetExecutor() queue.ExecFunc {
	return func(data interface{}) error {
		project, ok := data.(*configgithub.GitHubProjectConfig)
		if !ok {
			return ge.Pin(ge.New("data is not *configgithub.GitHubProjectConfig type"))
		}

		commands := []GitCommand{
			{
				Cmd:  project.BashScripts.OnPull.Before,
				Args: []string{},
			},
			{
				Cmd:  conn.config.GitCmd,
				Args: []string{"-C", project.Dir, "stash"},
			},
			{
				Cmd:  conn.config.GitCmd,
				Args: []string{"-C", project.Dir, "fetch"},
			},
			{
				Cmd:  conn.config.GitCmd,
				Args: []string{"-C", project.Dir, "checkout", project.Branch},
			},
			{
				Cmd:  conn.config.GitCmd,
				Args: []string{"-C", project.Dir, "pull"},
			},
			{
				Cmd:  project.BashScripts.OnPull.After,
				Args: []string{},
			},
		}

		for _, command := range commands {
			err := conn.exec(project.Repository, command.Cmd, command.Args)
			if err != nil {
				return ge.Pin(err)
			}
		}

		return nil
	}
}
