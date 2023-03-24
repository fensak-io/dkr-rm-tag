package dkrrmtag

import (
	"strings"

	"github.com/distribution/distribution/v3/reference"
)

type DockerRef struct {
	Host  string
	Owner string
	Name  string
	Tag   string
}

// ParseDockerImgRef parses the given user provided ref into the constituent parts. This handles the implicit rules
// implemented by the docker CLI, such as  `rust` => `docker.io/library/rust`.
func ParseDockerImgRef(ref string) (DockerRef, error) {
	parsedRef, err := reference.ParseNormalizedNamed(ref)
	if err != nil {
		return DockerRef{}, err
	}

	tag := "latest"
	if tagged, ok := parsedRef.(reference.Tagged); ok {
		tag = tagged.Tag()
	}

	fullPath := reference.Path(parsedRef)
	splitPath := strings.Split(fullPath, "/")
	owner := splitPath[0]
	name := strings.Join(splitPath[1:], "/")

	return DockerRef{
		Host:  reference.Domain(parsedRef),
		Owner: owner,
		Name:  name,
		Tag:   tag,
	}, nil
}

func (r DockerRef) AsDeleteTagRequest() DeleteTagRequest {
	return DeleteTagRequest{
		RepositoryOwner: r.Owner,
		RepositoryName:  r.Name,
		ImgTag:          r.Tag,
	}
}
