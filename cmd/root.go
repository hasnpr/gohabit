// Package cmd contains the command line interface for the app.
//
// Example:
//
//	$ ./gohabit start --config config.yaml
package cmd

import (
	"fmt"

	"github.com/hasnpr/gohabit/internal/app"
	"github.com/spf13/cobra"
)

var (
	// Flag variables
	cfgPath    string
	showBanner bool

	// rootCMD represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:                "gohabit",
		Short:              "A sample app that can be used",
		Long:               `A simple app that can be used as a sample that has redis, mariadb and tracing and logging`,
		PersistentPreRun:   preRun,
		PersistentPostRunE: postRun,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file path")
	rootCmd.PersistentFlags().BoolVar(&showBanner, "show-banner", false, "show application banner")
}

func preRun(_ *cobra.Command, _ []string) {
	fmt.Print(app.Banner(showBanner))
}

func postRun(_ *cobra.Command, _ []string) error {
	return nil
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
