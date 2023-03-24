package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fensak-io/dkr-rm-tag/version"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the version of dkr-rm-tag",
		Long:  `version prints out the version of the dkr-rm-tag binary.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("dkr-rm-tag version %s\n", version.Version)
			return nil
		},
	}
)
