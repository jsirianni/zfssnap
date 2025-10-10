// Package main implements the zfssnap CLI.
package main

import (
	"github.com/jsirianni/zfssnap/internal/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	RunE: func(_ *cobra.Command, _ []string) error {
		appLogger.Info("version", zap.String("version", version.Version()))
		return nil
	},
}
