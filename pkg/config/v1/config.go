package config

type Config struct {
	Github  `yaml:"github,omitempty"`
	Execute `yaml:"execute,omitempty"`
}

type Github struct {
	Owner             string `yaml:"owner,omitempty"`
	Repo              string `yaml:"repo,omitempty"`
	BotName           string `yaml:"botName,omitempty"`
	AccessTokenEnvVar string `yaml:"accessTokenEnvVar,omitempty"`
	PullRequest       `yaml:"pullRequest,omitempty"`
}

type PullRequest struct {
	Labels []string `yaml:"labels,omitempty"`
}

type Execute struct {
	Setup   []Command
	Track   []Command
	Cleanup []Command
}

type Command struct {
	Name string `yaml:"name,omitempty"`
	Cmd  string `yaml:"command,omitempty"`
	Dir  string `yaml:"dir,omitempty"`
}
