package dkrrmtag_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	dkrrmtag "github.com/fensak-io/dkr-rm-tag"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/environment"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestGHCRDeleteTag(t *testing.T) {
	t.Parallel()
	environment.RequireEnvVar(t, githubTokenEnvVarName)
	token := os.Getenv(githubTokenEnvVarName)
	ctx := context.Background()

	tag1ID := random.UniqueId()
	tag1 := fmt.Sprintf("%s:%s", testGHCRRepo, tag1ID)
	tag2ID := random.UniqueId()
	tag2 := fmt.Sprintf("%s:%s", testGHCRRepo, tag2ID)
	buildTestImage(t, tag1, tag2)

	reg := dkrrmtag.NewGHCR(token)
	req := dkrrmtag.DeleteTagRequest{
		RepositoryOwner: testDockerRepoOwner,
		RepositoryName:  testDockerRepoName,
		ImgTag:          tag1ID,
	}
	require.NoError(t, reg.DeleteTag(ctx, req))

	// Remove the local images, and then repull to test if the image was deleted from the registry
	docker.DeleteImage(t, tag1, nil)
	docker.DeleteImage(t, tag2, nil)

	opts := docker.RunOptions{}
	_, runErr1 := docker.RunE(t, tag1, &opts)
	require.Error(t, runErr1)
	docker.Run(t, tag2, &opts)
}
