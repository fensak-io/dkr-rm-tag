package main

import (
	"github.com/illumitacit/gostd/logstd"
)

type cmdOpts struct {
	Logger *logstd.Logger `mapstructure:"logger"`

	Tag string `mapstructure:"tag"`

	DockerHubUsername string `mapstructure:"docker_hub_username"`
	DockerHubToken    string `mapstructure:"docker_hub_token"`

	GitHubToken string `mapstructure:"github_token"`
}
