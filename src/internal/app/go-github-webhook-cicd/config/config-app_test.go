package configapp

import (
	"errors"
	configserver "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/config"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	configsentry "go-github-webhook-cicd/src/internal/pkg/services/sentry/config"
	"reflect"
	"testing"

	"github.com/nobuenhombre/suikat/pkg/fico"
)

type testConfig struct {
	fileName    string
	fileContent string
	config      *Config
	err         error
}

func TestConfigLoad(t *testing.T) {
	test := &testConfig{
		fileName:    "config-app_test_load.yaml",
		fileContent: "",
		config: &Config{
			Hosts: HostsConfig{
				API: configserver.HTTPServerConfig{
					Host: "127.0.0.1",
					Port: "7575",
				},
			},
			GitHub: configgithub.GitHubConfig{
				GitCmd: "/usr/bin/git",
				Projects: []configgithub.GitHubProjectConfig{
					{
						APIRoute:   "go-github-webhook-cicd-test-deploy",
						Repository: "git@github.com:nobuenhombre/go-github-webhook-cicd-test-deploy.git",
						Branch:     "main",
						Secret:     "tu38ir**_demo_**3if23ff",
						Dir:        "/opt/test-deploy/",
						BashScripts: configgithub.BashScriptsConfig{
							OnPull: configgithub.BashScriptsOnEventConfig{
								Before: "/opt/go-github-webhook-cicd/configs/develop.local.ivan/github/projects/test-deploy/bash-scripts/pull/after-pull.sh",
								After:  "/opt/go-github-webhook-cicd/configs/develop.local.ivan/github/projects/test-deploy/bash-scripts/pull/before-pull.sh",
							},
						},
					},
				},
			},
			Sentry: configsentry.SentryConfig{
				DSN:         "https://ABC@SUB.ingest.sentry.io/XXX",
				Environment: "develop",
			},
		},
		err: nil,
	}

	cfg := new(Config)
	err := cfg.Load(test.fileName)

	if !(reflect.DeepEqual(cfg, test.config) && errors.Is(err, test.err)) {
		t.Errorf(
			"cfg.Load(%#v),\n Expected (cfg = %#v, err = %#v),\n Actual (cfg = %#v, err = %#v).\n",
			test.fileName, test.config, test.err, cfg, err,
		)
	}
}

func TestConfigSave(t *testing.T) {
	test := &testConfig{
		fileName: "config-app_test_save.yaml",
		fileContent: "" +
			"hosts:\n" +
			"    api:\n" +
			"        host: 127.0.0.1\n" +
			"        port: \"7575\"\n" +
			"github:\n" +
			"    git_cmd: /usr/bin/git\n" +
			"    projects:\n" +
			"        - api_route: go-github-webhook-cicd-test-deploy\n" +
			"          repository: git@github.com:nobuenhombre/go-github-webhook-cicd-test-deploy.git\n" +
			"          branch: main\n" +
			"          secret: tu38ir**_demo_**3if23ff\n" +
			"          project_dir: /opt/test-deploy/\n" +
			"          bash_scripts:\n" +
			"            on_pull:\n" +
			"                before: /opt/go-github-webhook-cicd/configs/develop.local.ivan/github/projects/test-deploy/bash-scripts/pull/after-pull.sh\n" +
			"                after: /opt/go-github-webhook-cicd/configs/develop.local.ivan/github/projects/test-deploy/bash-scripts/pull/before-pull.sh\n" +
			"sentry:\n" +
			"    dsn: https://ABC@SUB.ingest.sentry.io/XXX\n" +
			"    environment: develop\n",
		config: &Config{
			Hosts: HostsConfig{
				API: configserver.HTTPServerConfig{
					Host: "127.0.0.1",
					Port: "7575",
				},
			},
			GitHub: configgithub.GitHubConfig{
				GitCmd: "/usr/bin/git",
				Projects: []configgithub.GitHubProjectConfig{
					{
						APIRoute:   "go-github-webhook-cicd-test-deploy",
						Repository: "git@github.com:nobuenhombre/go-github-webhook-cicd-test-deploy.git",
						Branch:     "main",
						Secret:     "tu38ir**_demo_**3if23ff",
						Dir:        "/opt/test-deploy/",
						BashScripts: configgithub.BashScriptsConfig{
							OnPull: configgithub.BashScriptsOnEventConfig{
								Before: "/opt/go-github-webhook-cicd/configs/develop.local.ivan/github/projects/test-deploy/bash-scripts/pull/after-pull.sh",
								After:  "/opt/go-github-webhook-cicd/configs/develop.local.ivan/github/projects/test-deploy/bash-scripts/pull/before-pull.sh",
							},
						},
					},
				},
			},
			Sentry: configsentry.SentryConfig{
				DSN:         "https://ABC@SUB.ingest.sentry.io/XXX",
				Environment: "develop",
			},
		},
		err: nil,
	}

	cfg := test.config
	err := cfg.Save(test.fileName)

	txtConfigFile := fico.TxtFile(test.fileName)
	fileContent, errReadFile := txtConfigFile.Read()

	if errReadFile != nil {
		t.Errorf(
			"txtConfigFile.Read error %#v",
			errReadFile,
		)
	}

	if !(reflect.DeepEqual(fileContent, test.fileContent) && errors.Is(err, test.err)) {
		t.Errorf(
			"cfg.Save(%#v),\n Expected (fileContent = %#v, err = %#v),\n Actual (fileContent = %#v, err = %#v).\n",
			test.fileName, test.fileContent, test.err, fileContent, err,
		)
	}
}
