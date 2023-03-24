# dkr-rm-tag

`dkr-rm-tag` is a Go based CLI that can be used to remove a specific Docker image tag from a target registry.

Ideally, you would be able to use the `docker` CLI to implement this functionality, but as of March 24th, 2023, the
`docker` CLI only offers the ability to remove image tags locally and not in the remote registry. Instead, each registry
offers a proprietary API endpoint to implement this feature.

This CLI is meant to provide a single unified interface for removing image tags from various underlying registries.


## Quick start

TODO


## Supported registries

- [Docker Hub](https://hub.docker.com/)
- [GitHub Packages (GHCR)](https://github.com/features/packages)
