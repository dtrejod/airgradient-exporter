package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dtrejod/airgradient-exporter/internal/ilog"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// NewAirGradient creates a new collector for the AirGradient local server API.
// https://github.com/airgradienthq/arduino/blob/master/docs/local-server.md#local-server-api
func NewAirGradient(ctx context.Context, endpoint string) (prometheus.Collector, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse airgradient endpoint into url: %w", err)
	}

	return &airgradientCollector{
		ctx:      ctx,
		client:   &http.Client{},
		endpoint: e,
		deviceInfoDesc: prometheus.NewDesc(
			"airgradient_device_info",
			"Device information",
			[]string{"serialno", "firmware", "model", "ledmode"},
			nil,
		),
		wifiDesc: prometheus.NewDesc(
			"airgradient_wifi",
			"WiFi signal strength",
			nil,
			nil,
		),
		pm01Desc: prometheus.NewDesc(
			"airgradient_pm01",
			"PM1 in ug/m3",
			nil,
			nil,
		),
		pm02Desc: prometheus.NewDesc(
			"airgradient_pm02",
			"PM2.5 in ug/m3",
			nil,
			nil,
		),
		pm10Desc: prometheus.NewDesc(
			"airgradient_pm10",
			"PM10 in ug/m3",
			nil,
			nil,
		),
		pm02CompensatedDesc: prometheus.NewDesc(
			"airgradient_pm02_compensated",
			"PM2.5 in ug/m3 with correction applied",
			nil,
			nil,
		),
		rco2Desc: prometheus.NewDesc(
			"airgradient_rco2",
			"CO2 in ppm",
			nil,
			nil,
		),
		pm003CountDesc: prometheus.NewDesc(
			"airgradient_pm003_count",
			"Particle count per dL",
			nil,
			nil,
		),
		atmpDesc: prometheus.NewDesc(
			"airgradient_atmp",
			"Temperature in Degrees Celsius",
			nil,
			nil,
		),
		atmpCompensatedDesc: prometheus.NewDesc(
			"airgradient_atmp_compensated",
			"Temperature in Degrees Celsius with correction applied",
			nil,
			nil,
		),
		rhumDesc: prometheus.NewDesc(
			"airgradient_rhum",
			"Relative Humidity",
			nil,
			nil,
		),
		rhumCompensatedDesc: prometheus.NewDesc(
			"airgradient_rhum_compensated",
			"Relative Humidity with correction applied",
			nil,
			nil,
		),
		tvocIndexDesc: prometheus.NewDesc(
			"airgradient_tvoc_index",
			"Senisiron VOC Index",
			nil,
			nil,
		),
		tvocRawDesc: prometheus.NewDesc(
			"airgradient_tvoc_raw",
			"VOC raw value",
			nil,
			nil,
		),
		noxIndexDesc: prometheus.NewDesc(
			"airgradient_nox_index",
			"Senisirion NOx Index",
			nil,
			nil,
		),
		noxRawDesc: prometheus.NewDesc(
			"airgradient_nox_raw",
			"NOx raw value",
			nil,
			nil,
		),
		bootDesc: prometheus.NewDesc(
			"airgradient_boot",
			"Counts every measurement cycle. Low boot counts indicate restarts",
			nil,
			nil,
		),
		bootCountDesc: prometheus.NewDesc(
			"airgradient_boot_count",
			"Same as boot property. Required for Home Assistant compatability. Will be depreciated",
			nil,
			nil,
		),
	}, nil
}

type airgradientCollector struct {
	ctx      context.Context
	client   *http.Client
	endpoint *url.URL

	deviceInfoDesc      *prometheus.Desc
	wifiDesc            *prometheus.Desc
	pm01Desc            *prometheus.Desc
	pm02Desc            *prometheus.Desc
	pm10Desc            *prometheus.Desc
	pm02CompensatedDesc *prometheus.Desc
	rco2Desc            *prometheus.Desc
	pm003CountDesc      *prometheus.Desc
	atmpDesc            *prometheus.Desc
	atmpCompensatedDesc *prometheus.Desc
	rhumDesc            *prometheus.Desc
	rhumCompensatedDesc *prometheus.Desc
	tvocIndexDesc       *prometheus.Desc
	tvocRawDesc         *prometheus.Desc
	noxIndexDesc        *prometheus.Desc
	noxRawDesc          *prometheus.Desc
	bootDesc            *prometheus.Desc
	bootCountDesc       *prometheus.Desc
}

func (c *airgradientCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.deviceInfoDesc
	ch <- c.wifiDesc
	ch <- c.pm01Desc
	ch <- c.pm02Desc
	ch <- c.pm10Desc
	ch <- c.pm02CompensatedDesc
	ch <- c.rco2Desc
	ch <- c.pm003CountDesc
	ch <- c.atmpDesc
	ch <- c.atmpCompensatedDesc
	ch <- c.rhumDesc
	ch <- c.rhumCompensatedDesc
	ch <- c.tvocIndexDesc
	ch <- c.tvocRawDesc
	ch <- c.noxIndexDesc
	ch <- c.noxRawDesc
	ch <- c.bootDesc
	ch <- c.bootCountDesc
}

func (c *airgradientCollector) Collect(ch chan<- prometheus.Metric) {
	m, err := c.getMeasures(c.ctx)
	if err != nil {
		ilog.FromContext(c.ctx).Error("Failed to get measures.", zap.Error(err))
		return
	}

	ch <- prometheus.MustNewConstMetric(c.deviceInfoDesc, prometheus.GaugeValue, 1, m.SerialNo, m.Firmware, m.Model, m.LEDMode)
	ch <- prometheus.MustNewConstMetric(c.wifiDesc, prometheus.GaugeValue, float64(m.Wifi))
	ch <- prometheus.MustNewConstMetric(c.pm01Desc, prometheus.GaugeValue, float64(m.PM01))
	ch <- prometheus.MustNewConstMetric(c.pm02Desc, prometheus.GaugeValue, float64(m.PM02))
	ch <- prometheus.MustNewConstMetric(c.pm10Desc, prometheus.GaugeValue, float64(m.PM10))
	ch <- prometheus.MustNewConstMetric(c.pm02CompensatedDesc, prometheus.GaugeValue, float64(m.PM02Compensated))
	ch <- prometheus.MustNewConstMetric(c.rco2Desc, prometheus.GaugeValue, float64(m.RCO2))
	ch <- prometheus.MustNewConstMetric(c.pm003CountDesc, prometheus.GaugeValue, float64(m.PM003Count))
	ch <- prometheus.MustNewConstMetric(c.atmpDesc, prometheus.GaugeValue, m.ATMP)
	ch <- prometheus.MustNewConstMetric(c.atmpCompensatedDesc, prometheus.GaugeValue, m.ATMPCompensated)
	ch <- prometheus.MustNewConstMetric(c.rhumDesc, prometheus.GaugeValue, float64(m.RHUM))
	ch <- prometheus.MustNewConstMetric(c.rhumCompensatedDesc, prometheus.GaugeValue, float64(m.RHUMCompensated))
	ch <- prometheus.MustNewConstMetric(c.tvocIndexDesc, prometheus.GaugeValue, float64(m.TVOCIndex))
	ch <- prometheus.MustNewConstMetric(c.tvocRawDesc, prometheus.GaugeValue, float64(m.TVOCRaw))
	ch <- prometheus.MustNewConstMetric(c.noxIndexDesc, prometheus.GaugeValue, float64(m.NOXIndex))
	ch <- prometheus.MustNewConstMetric(c.noxRawDesc, prometheus.GaugeValue, float64(m.NOXRaw))
	ch <- prometheus.MustNewConstMetric(c.bootDesc, prometheus.GaugeValue, float64(m.Boot))
	ch <- prometheus.MustNewConstMetric(c.bootCountDesc, prometheus.GaugeValue, float64(m.BootCount))
}

func (c *airgradientCollector) getMeasures(ctx context.Context) (*measures, error) {
	ilog.FromContext(ctx).Debug("Getting measures from airgradient.")
	req, err := http.NewRequestWithContext(ctx, "GET", c.endpoint.JoinPath(measuresPath).String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m measures
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	ilog.FromContext(ctx).Debug("Got measures from airgradient.", zap.Any("measures", m))
	return &m, nil
}
