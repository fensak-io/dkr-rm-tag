package dkrrmtag_test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
)

const (
	dkrHubTokenEnvVarName = "DOCKER_HUB_TOKEN"
	testDockerRepoOwner   = "masoukishi"
	testDockerRepoName    = "dkr-rm-tag-test"

	githubTokenEnvVarName = "GITHUB_TOKEN"
	testGHCRRepoOwner     = "fensak-io"
	testGHCRRepoName      = "dkr-rm-tag-test"
)

var (
	testDockerRepo = fmt.Sprintf("%s/%s", testDockerRepoOwner, testDockerRepoName)
	testGHCRRepo   = fmt.Sprintf("ghcr.io/%s/%s", testGHCRRepoOwner, testGHCRRepoName)
)

func buildTestImage(t *testing.T, tags ...string) {
	opts := docker.BuildOptions{
		Tags:           tags,
		EnableBuildKit: true,
		Push:           true,
	}
	docker.Build(t, "./fixtures", &opts)
}
