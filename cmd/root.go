package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/abhishek/dmon-cli/internal/logger"
	"github.com/abhishek/dmon-cli/internal/service"
	"github.com/abhishek/dmon-cli/internal/xrandr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	log     *logrus.Logger
	svc     *service.DisplayService
)

var rootCmd = &cobra.Command{
	Use:   "dmon",
	Short: "Display Monitor - manage your displays with ease",
	Long: `dmon is a CLI tool for managing display configurations on Linux.
It provides a simple interface to xrandr for common display management tasks.`,
	Version: "1.0.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		log, err = logger.New(verbose)
		if err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		backend := xrandr.NewBackend(log)
		svc = service.New(backend, log)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed output and xrandr commands")
}

func getContext() context.Context {
	return context.Background()
}
