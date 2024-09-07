package collector

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

// NewAirGradient creates a new collector for the AirGradient local server API.
// https://github.com/airgradienthq/arduino/blob/master/docs/local-server.md#local-server-api
func NewAirGradient(ctx context.Context, endpoint string) (prometheus.Collector, error) {
}

type airgradientCollector struct {
	ctx context.Context

	//billUsageKwhDesc            *prometheus.Desc
	//billUSDServiceChargeDesc    *prometheus.Desc
	//currentUsageKwhDesc         *prometheus.Desc
	//currentUSDServiceChargeDesc *prometheus.Desc
}

func (c *airgradientCollector) Describe(ch chan<- *prometheus.Desc) {
	return
}

func (c *airgradientCollector) Collect(ch chan<- prometheus.Metric) {
	return
}
