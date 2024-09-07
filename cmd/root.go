package cmd

import (
	"context"

	"github.com/dtrejod/airgradient-exporter/internal/ilog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ctx     context.Context
	debug   bool
	rootCmd = &cobra.Command{
		Use:               "airgradient",
		Short:             "A Prometheus exporter for AirGradient local server",
		PersistentPreRunE: initLoggers,
	}
)

func initLoggers(_ *cobra.Command, _ []string) error {
	ctx = context.Background()
	var err error
	logConfig := zap.NewProductionConfig()
	if debug {
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := logConfig.Build()
	if err != nil {
		return err
	}
	ctx = ilog.WithLogger(ctx, logger)
	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "run in debug mode")
}

// Execute runs the root command tree
func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		ilog.FromContext(ctx).Error("Command error", zap.Error(err))
		return 1
	}
	return 0
}
