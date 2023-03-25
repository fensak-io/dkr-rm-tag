package dkrrmtag

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
)

// dkrHubRegistry represents an HTTP client that can communicate with the Docker Hub registry API. Note that this is
// different from the docker engine API, and provides access to additional endpoints not accessible from the docker CLI.
type dkrHubRegistry struct {
	clt   *resty.Client
	token string
}

var _ (Registry) = (*dkrHubRegistry)(nil)

// loginRequest represents the POST body for the login request to Docker Hub.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// loginResponse represents the response object from a successful login request to the Docker Hub API.
type loginResponse struct {
	Token string `json:"token"`
}

// NewDkrHubRegistry returns a handle to a new dkrHubRegistry object which can be used to make requests to the registry.
func NewDkrHubRegistry(ctx context.Context, username, token string) (dkrHubRegistry, error) {
	clt := resty.New()

	// Use the auth endpoint to get a token to authenticate to Docker Hub API.
	// See https://docs.docker.com/docker-hub/api/latest/#tag/authentication for more information.
	loginURL := dkrHubAPIURL("users/login")
	req := loginRequest{Username: username, Password: token}
	httpResp, err := clt.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&loginResponse{}).
		Post(loginURL)
	if err != nil {
		return dkrHubRegistry{}, err
	}

	// NOTE: we don't do a type check here as we trust resty to do the right thing.
	resp := httpResp.Result().(*loginResponse)
	return dkrHubRegistry{clt: clt, token: resp.Token}, nil
}

// DeleteTag uses an undocumented API on Docker Hub to delete the image tag from the remote repository. Although this is
// not documented, empirically it seems to work. However, it is expected that this could break at any time so caution is
// needed.
func (r dkrHubRegistry) DeleteTag(ctx context.Context, req DeleteTagRequest) error {
	urlpath := path.Join(
		"repositories",
		req.RepositoryOwner,
		req.RepositoryName,
		"tags",
		req.ImgTag,
	)
	resp, err := r.clt.R().
		SetContext(ctx).
		SetAuthToken(r.token).
		Delete(dkrHubAPIURL(urlpath))
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf(
			"Received error status code making delete call to Docker: code %d, body %s",
			resp.StatusCode(), string(resp.Body()),
		)
	}
	return nil
}

// dkrHubAPIURL returns the full URL to the docker hub API given the relative path. This is hard coded to:
// - hub.docker.com
// - v2 API
func dkrHubAPIURL(relPath string) string {
	u := url.URL{
		Scheme: "https",
		Host:   "hub.docker.com",
		Path:   path.Join("v2", relPath),
	}
	return u.String()
}
