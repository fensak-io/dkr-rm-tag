package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	dkrrmtag "github.com/fensak-io/dkr-rm-tag"
	"github.com/fensak-io/gostd/clistd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	pflags := rootCmd.PersistentFlags()

	// Configuration options for the logger
	pflags.String("loglevel", "", "Logging level. Valid options: debug, info, warn, error, panic, fatal")
	clistd.MustBindPFlag("logger.level", pflags.Lookup("loglevel"))

	pflags.String("logencoding", "console", "Log message encoding. Valid options: json, console")
	clistd.MustBindPFlag("logger.encoding", pflags.Lookup("logencoding"))

	// Main routine options
	pflags.String("tag", "", "The full remote docker image tag to remove. This should include the registry URL and repository information (e.g., ghcr.io/fensak-io/dkr-rm-tag:sometag).")
	clistd.MustBindPFlag("tag", pflags.Lookup("tag"))

	pflags.String("docker-hub-username", "", "The username to use for authenticating to Docker Hub.")
	clistd.MustBindPFlag("docker_hub_username", pflags.Lookup("docker-hub-username"))

	pflags.String("docker-hub-token", "", "The token to use for authenticating to Docker Hub. Recommended to be passed in by environment variable (DOCKER_HUB_TOKEN).")
	clistd.MustBindPFlag("docker_hub_token", pflags.Lookup("docker-hub-token"))
	viper.BindEnv("docker_hub_token", "DOCKER_HUB_TOKEN")

	pflags.String("github-token", "", "The token to use for authenticating to GitHub Container Registry. Recommended to be passed in by environment variable (GITHUB_TOKEN).")
	clistd.MustBindPFlag("github_token", pflags.Lookup("github-token"))
	viper.BindEnv("github_token", "GITHUB_TOKEN")
}

var rootCmd = &cobra.Command{
	Use:   "dkr-rm-tag",
	Short: "dkr-rm-tag is a CLI for deleting an image tag from a remote Docker registry",
	Long:  "dkr-rm-tag is a CLI for deleting an image tag from a remote Docker registry.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var opts cmdOpts
		if err := viper.Unmarshal(&opts); err != nil {
			return err
		}

		// Validate inputs
		if opts.Tag == "" {
			return errors.New("--tag is required")
		}

		// Arbitrary 5 minute timeout for all operations
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		parsed, err := dkrrmtag.ParseDockerImgRef(opts.Tag)
		if err != nil {
			return err
		}

		var reg dkrrmtag.Registry
		switch parsed.Host {
		case "docker.io":
			dkrReg, err := dkrrmtag.NewDkrHubRegistry(ctx, opts.DockerHubUsername, opts.DockerHubToken)
			if err != nil {
				return err
			}
			reg = dkrReg
		case "ghcr.io":
			reg = dkrrmtag.NewGHCR(opts.GitHubToken)
		default:
			return fmt.Errorf("registry hosted at %s is not supported by dkr-rm-tag", parsed.Host)
		}

		return reg.DeleteTag(ctx, parsed.AsDeleteTagRequest())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
