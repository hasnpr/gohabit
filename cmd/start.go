// Package cmd contains the command line interface for the app.
//
// Example:
//
//	$ ./gohabit start --config config.yaml
package cmd

import (
	"github.com/hasnpr/gohabit/internal/app"
	"github.com/spf13/cobra"
)

// startCMD represents the start command of the app.
var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "A sample app that can be used",
		Long:  `A simple app that can be used as a sample that has redis, mariadb and tracing and logging`,
		Run:   start,
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(_ *cobra.Command, _ []string) {
	app.Start(cmdLogger)
}
