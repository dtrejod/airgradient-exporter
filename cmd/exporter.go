package cmd

import (
	"net/http"
	"os"

	"github.com/dtrejod/airgradient-exporter/internal/collector"
	"github.com/dtrejod/airgradient-exporter/internal/ilog"
	"github.com/dtrejod/airgradient-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	ilog.FromContext(ctx).Info("Starting airgradient-exporter...", zap.String("version", version.Version()))

	if endpoint == "" {
		ilog.FromContext(ctx).Fatal("Missing required '--endpoint' arguement.")
		os.Exit(1)
	}

	airgradientCollector, err := collector.NewAirGradient(ctx, endpoint)
	if err != nil {
		ilog.FromContext(ctx).Fatal("Failed to create airgradient-exporter.", zap.Error(err))
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
	exporterCmd.Flags().StringVar(&endpoint, endpointFlag, "", "AirGradient local-server endpoint. (e.g http://airgradient_<serial-number>.local)")
	if err := viper.BindPFlag(endpointFlag, exporterCmd.Flags().Lookup(endpointFlag)); err != nil {
		panic(err)
	}
	if err := viper.BindEnv(endpointFlag, "ENDPOINT"); err != nil {
		panic(err)
	}
	endpoint = viper.GetString(endpointFlag)

	exporterCmd.Flags().StringVar(&listenAddr, listenAddrFlag, ":9091", "HTTP port to listen on.")
	if err := viper.BindPFlag(listenAddrFlag, exporterCmd.Flags().Lookup(listenAddrFlag)); err != nil {
		panic(err)
	}
	if err := viper.BindEnv(listenAddrFlag, "LISTEN_ADDRESS"); err != nil {
		panic(err)
	}
	listenAddr = viper.GetString(listenAddrFlag)

	rootCmd.AddCommand(exporterCmd)
}
