package configgithub

type BashScriptsOnEventConfig struct {
	Before string `yaml:"before,omitempty"`
	After  string `yaml:"after,omitempty"`
}

type BashScriptsConfig struct {
	OnPull BashScriptsOnEventConfig `yaml:"on_pull,omitempty"`
}

type GitHubProjectConfig struct {
	APIRoute    string            `yaml:"api_route,omitempty"`
	Repository  string            `yaml:"repository,omitempty"`
	Branch      string            `yaml:"branch,omitempty"`
	Token       string            `yaml:"token,omitempty"`
	Dir         string            `yaml:"project_dir,omitempty"`
	BashScripts BashScriptsConfig `yaml:"bash_scripts,omitempty"`
}

type GitHubConfig struct {
	GitCmd   string                `yaml:"git_cmd,omitempty"`
	Projects []GitHubProjectConfig `yaml:"projects,omitempty"`
}
