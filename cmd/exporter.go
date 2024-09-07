package cmd

import (
	"net/http"
	"os"

	"github.com/dtrejod/airgradient-exporter/internal/collector"
	"github.com/dtrejod/airgradient-exporter/internal/ilog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	listenAddrFlag = "listen-address"
	endpointFlag   = "endpoint"
)

var (
	listenAddr string
	endpoint   string
)

var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "Run the Exporter",
	Run:   exporterRunFunc,
}

func exporterRunFunc(cmd *cobra.Command, args []string) {
	airgradientCollector, err := collector.NewAirGradient(ctx, endpoint)
	if err != nil {
		ilog.FromContext(ctx).Fatal("Failed to create airgradient-collector.", zap.Error(err))
		os.Exit(1)
	}
	if err := prometheus.Register(airgradientCollector); err != nil {
		ilog.FromContext(ctx).Fatal("Failed to register collector.", zap.Error(err))
		os.Exit(1)
	}
	http.Handle("/metrics", promhttp.Handler())

	ilog.FromContext(ctx).Info("Starting server", zap.String("addr", listenAddr))
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		ilog.FromContext(ctx).Fatal("Failed to start exporter server", zap.Error(err))
		os.Exit(1)
	}
	ilog.FromContext(ctx).Info("Exporter server stopped.")
}

func init() {
	exporterCmd.Flags().StringVar(&listenAddr, listenAddrFlag, ":9091", "HTTP port to listen on.")
	exporterCmd.Flags().StringVar(&endpoint, endpointFlag, "", "AirGradient localserver endpoint. (e.g http://airgradient_<serial-number>.local)")
	_ = exporterCmd.MarkFlagRequired(endpointFlag)
	rootCmd.AddCommand(exporterCmd)
}
