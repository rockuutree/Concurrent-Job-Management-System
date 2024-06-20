package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	version "uw.cse374.go/scheduler"
	"uw.cse374.go/scheduler/internal/config"
	"uw.cse374.go/scheduler/internal/scheduler"
)

// errInvalidArgs is used when the user doesn't specify a single argument.
var errInvalidArgs = errors.New("exactly one argument is required for this command")

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

// run is the primary entrypoint to the scheduler program. This sets up
// the CLI interface, and dispatches the sub-command that was specified on
// the command line (e.g. ./scheduler run).
func run() error {
	rootCmd := &cobra.Command{
		Use:     "scheduler",
		Short:   "Job Scheduler",
		Version: version.Version,
	}

	// Define the runCmd with its RunE function
	runCmd := &cobra.Command{
		Use:   "run <file>",
		Short: "Run the configured jobs",
		// Changed RunE to use the cmd.Flags()
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errInvalidArgs
			}

			//  Changed to use the cmd.Flags() to get the debug flag
			debug, _ := cmd.Flags().GetBool("debug")

			configFile := args[0]
			scheduler, err := newScheduler(configFile)
			if err != nil {
				return err
			}

			return scheduler.Run(os.Stdout, debug)
		},
	}

	// Add debug flag to the runCmd
	runCmd.Flags().BoolP("debug", "d", false, "enable debug logging")

	// Define the graphCmd with its RunE function
	graphCmd := &cobra.Command{
		Use:   "graph <file>",
		Short: "Write a DOT-graph for the configured jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errInvalidArgs
			}

			configFile := args[0]
			scheduler, err := newScheduler(configFile)
			if err != nil {
				return err
			}

			return scheduler.Graph(os.Stdout)
		},
	}

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(graphCmd)

	// Rewrite the --version output so that we only
	// produce a semantic version.
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	// Hide and/or disable the commands we don't care about.
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	// Silence usage and error messages on error.
	rootCmd.SilenceUsage = true

	return rootCmd.Execute()
}

func newScheduler(configFile string) (*scheduler.Scheduler, error) {
	config, err := config.ReadConfig(configFile)
	if err != nil {
		return nil, err
	}

	return scheduler.New(config), nil
}
