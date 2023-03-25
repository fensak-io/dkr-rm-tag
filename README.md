<h1 align="center">dkr-rm-tag</h1>

<p align="center">
  <a href="https://github.com/fensak-io/dkr-rm-tag/blob/main/LICENSE">
    <img alt="LICENSE" src="https://img.shields.io/github/license/fensak-io/dkr-rm-tag?style=for-the-badge">
  </a>
  <a href="https://github.com/fensak-io/dkr-rm-tag/actions/workflows/lint-and-test.yml?query=branch%3Amain">
    <img alt="main branch CI" src="https://img.shields.io/github/actions/workflow/status/fensak-io/dkr-rm-tag/lint-and-test.yml?branch=main&logo=github&label=CI&style=for-the-badge">
  </a>
  <a href="https://github.com/fensak-io/dkr-rm-tag/releases/latest">
    <img alt="latest release" src="https://img.shields.io/github/v/release/fensak-io/dkr-rm-tag?style=for-the-badge">
  </a>
</p>

`dkr-rm-tag` is a Go based CLI that can be used to remove a specific Docker image tag from a target registry.

Ideally, you would be able to use the `docker` CLI to implement this functionality, but as of March 24th, 2023, the
`docker` CLI only offers the ability to remove image tags locally and not in the remote registry. Instead, each registry
offers a proprietary API endpoint to implement this feature.

This CLI is meant to provide a single unified interface for removing image tags from various underlying registries.


## Quick start

Use one of the following to download the `dkr-rm-tag` CLI:

- Downloading one of the pre-compiled binaries from the [releases page](/releases).
- Building from source using `go`:

      go install github.com/fensak-io/dkr-rm-tag/cmd/dkr-rm-tag@latest

Once you have the CLI installed, you can run the command to remove a tag from one of the supported registries. For
example, to remove the `user/myrepo:some-tag` tag from Docker Hub:

```
  export DOCKER_HUB_TOKEN='....your personal access token...'
dkr-rm-tag --tag 'user/myrepo:some-tag' --docker-hub-username user
```

Alternatively, to remove from `ghcr.io`:

```
  export GITHUB_TOKEN='...your personal access token...'
dkr-rm-tag --tag 'ghcr.io/user/myrepo:some-tag'
```

Refer to `dkr-rm-tag --help` for all the available options by the command.


## Supported registries

- [Docker Hub](https://hub.docker.com/)
- [GitHub Packages (GHCR)](https://github.com/features/packages)
