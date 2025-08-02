// Package cmd contains the command line interface for the app.
//
// Example:
//
//	$ ./rasoul start --config config.yaml
package cmd

import (
	"fmt"

	"github.com/hasnpr/gohabit/internal/app"
	"github.com/hasnpr/gohabit/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	cmdLogger *logger.Logger

	// Flag variables
	cfgPath  string
	noBanner bool

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
	rootCmd.PersistentFlags().BoolVar(&noBanner, "no-banner", false, "hide application banner")

	cmdLogger = logger.NewDefault()
}

func preRun(_ *cobra.Command, _ []string) {
	fmt.Print(app.Banner(noBanner))
}

func postRun(_ *cobra.Command, _ []string) error {
	return cmdLogger.Close()
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var rootCMD = &cobra.Command{
	Use:   "gohabit",
	Short: "GoHabit service",
}
