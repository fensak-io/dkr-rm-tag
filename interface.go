package dkrrmtag

import "context"

// Registry provides an interface into a docker registry that supports deleting image tags.
type Registry interface {
	// DeleteTag deletes a given image tag from the underlying docker registry.
	DeleteTag(context.Context, DeleteTagRequest) error
}

// DeleteTagRequest specifies the image tag from a given repository that should be deleted from the registry.
type DeleteTagRequest struct {
	RepositoryOwner string
	RepositoryName  string
	ImgTag          string
}
