package dkrrmtag

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

// ghcr represents a GraphQL client that can communicate with the GitHub packages API for managing docker images.
type ghcr struct {
	clt *github.Client
}

var _ (Registry) = (*ghcr)(nil)

// NewGHCR returns a handle to a new ghcr object which can be used to make requests to the registry (using the GraphQL
// API).
func NewGHCR(token string) ghcr {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := github.NewClient(httpClient)

	return ghcr{clt: client}
}

// DeleteTag uses the GraphQL API to remove an image tag from the GHCR docker registry. This makes two calls: once to
// get the corresponding object ID of the image tag, and then the delete request using the ID.
func (r ghcr) DeleteTag(ctx context.Context, req DeleteTagRequest) error {
	ownerType, err := r.determineOwnerType(ctx, req.RepositoryOwner)
	if err != nil {
		return err
	}

	var v *github.PackageVersion
	var findErr error
	switch ownerType {
	case ownerTypeOrg:
		v, findErr = r.findPackageVersionWithTag(ctx, req, r.clt.Organizations.PackageGetAllVersions)
	case ownerTypeUser:
		v, findErr = r.findPackageVersionWithTag(ctx, req, r.clt.Users.PackageGetAllVersions)
	default:
		return errors.New("impossible condition")
	}
	if findErr != nil {
		return findErr
	}

	versionTags := v.Metadata.Container.Tags

	// If this is the only tag in the version, then delete the package and return. Otherwise, go through a special routine
	// to ensure the other tags are preserved.
	if len(versionTags) == 1 {
		return r.deletePackageVersion(ctx, ownerType, *v.ID, req)
	}

	// GHCR doesn't have an API for deleting just a single image tag. Their API offers a way to delete a package version,
	// but invoking that would cause the API to delete the whole image, which is problematic if there are multiple tags
	// referencing that version.
	// To handle this, we will do the following:
	// - Use docker to pull down the list of tags we need to keep.
	// - Delete the package version.
	// - Push the tags remaining tags using docker to restore the package version, without the tag we want to remove.
	// Note that this could lead to brief moments of downtime since the other tags would be missing for a period of time.
	tags := make([]string, 0, len(versionTags)-1)
	for _, t := range versionTags {
		tags = append(tags, fmt.Sprintf("ghcr.io/%s/%s:%s", req.RepositoryOwner, req.RepositoryName, t))
	}
	for _, t := range tags {
		if err := dockerPull(t); err != nil {
			return err
		}
	}
	if err := r.deletePackageVersion(ctx, ownerType, *v.ID, req); err != nil {
		return err
	}
	for _, t := range tags {
		if err := dockerPush(t); err != nil {
			return err
		}
	}
	return nil
}

// findPackageVersionWithTag pages through all versions of a repository looking for the specific tag.
func (r ghcr) findPackageVersionWithTag(
	ctx context.Context, req DeleteTagRequest,
	pkgVersionListFunc func(context.Context, string, string, string, *github.PackageListOptions) ([]*github.PackageVersion, *github.Response, error),
) (*github.PackageVersion, error) {
	opts := &github.PackageListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	const maxPages = 1000
	var i int
	for i = 0; i < maxPages; i++ {
		versions, resp, err := pkgVersionListFunc(ctx, req.RepositoryOwner, "container", req.RepositoryName, opts)
		if err != nil {
			return nil, err
		}
		for _, v := range versions {
			if v.Metadata != nil && v.Metadata.Container != nil && slices.Contains(v.Metadata.Container.Tags, req.ImgTag) {
				return v, nil
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	if i >= 1000 {
		return nil, errors.New("reached max package version page limit")
	}
	return nil, fmt.Errorf(
		"GHCR container repo %s/%s does not have tag %s",
		req.RepositoryOwner, req.RepositoryName, req.ImgTag,
	)
}

// deletePackageVersion uses the github API to delete the requested package version identified by packageVersionID from
// the requested repository.
func (r ghcr) deletePackageVersion(ctx context.Context, ot ownerType, packageVersionID int64, req DeleteTagRequest) error {
	var fn func(context.Context, string, string, string, int64) (*github.Response, error)
	switch ot {
	case ownerTypeOrg:
		fn = r.clt.Organizations.PackageDeleteVersion
	case ownerTypeUser:
		fn = r.clt.Users.PackageDeleteVersion
	default:
		return errors.New("impossible condition")
	}
	_, err := fn(ctx, req.RepositoryOwner, "container", req.RepositoryName, packageVersionID)
	return err
}

// determineOwnerType determines if the given slug refers to a github org, or user.
func (r ghcr) determineOwnerType(ctx context.Context, slug string) (ownerType, error) {
	if _, resp, _ := r.clt.Organizations.Get(ctx, slug); resp.StatusCode == http.StatusOK {
		return ownerTypeOrg, nil
	}

	if _, resp, _ := r.clt.Users.Get(ctx, slug); resp.StatusCode == http.StatusOK {
		return ownerTypeUser, nil
	}

	return ownerTypeUnknown, fmt.Errorf("The github slug %s refers to neither an org or a user.", slug)
}

type ownerType int

const (
	ownerTypeUnknown ownerType = iota
	ownerTypeOrg
	ownerTypeUser
)
