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
func NewAirGradient(ctx context.Context, endpoint string) (*AirgradientCollector, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not parse airgradient endpoint into url: %w", err)
	}

	return &AirgradientCollector{
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
			[]string{"serialno"},
			nil,
		),
		pm01Desc: prometheus.NewDesc(
			"airgradient_pm01",
			"PM1 in ug/m3",
			[]string{"serialno"},
			nil,
		),
		pm01StandardDesc: prometheus.NewDesc(
			"airgradient_pm01_standard",
			"PM1 in ug/m3 (standard particle)",
			[]string{"serialno"},
			nil,
		),
		pm01CountDesc: prometheus.NewDesc(
			"airgradient_pm01_count",
			"Particle count 1um per dL",
			[]string{"serialno"},
			nil,
		),
		pm02Desc: prometheus.NewDesc(
			"airgradient_pm02",
			"PM2.5 in ug/m3",
			[]string{"serialno"},
			nil,
		),
		pm02StandardDesc: prometheus.NewDesc(
			"airgradient_pm02_standard",
			"PM2.5 in ug/m3 (standard particle)",
			[]string{"serialno"},
			nil,
		),
		pm02CountDesc: prometheus.NewDesc(
			"airgradient_pm02_count",
			"Particle count 2.5um per dL",
			[]string{"serialno"},
			nil,
		),
		pm10Desc: prometheus.NewDesc(
			"airgradient_pm10",
			"PM10 in ug/m3",
			[]string{"serialno"},
			nil,
		),
		pm10StandardDesc: prometheus.NewDesc(
			"airgradient_pm10_standard",
			"PM10 in ug/m3 (standard particle)",
			[]string{"serialno"},
			nil,
		),
		pm02CompensatedDesc: prometheus.NewDesc(
			"airgradient_pm02_compensated",
			"PM2.5 in ug/m3 with correction applied",
			[]string{"serialno"},
			nil,
		),
		rco2Desc: prometheus.NewDesc(
			"airgradient_rco2",
			"CO2 in ppm",
			[]string{"serialno"},
			nil,
		),
		pm003CountDesc: prometheus.NewDesc(
			"airgradient_pm003_count",
			"Particle count 0.3um per dL",
			[]string{"serialno"},
			nil,
		),
		pm005CountDesc: prometheus.NewDesc(
			"airgradient_pm005_count",
			"Particle count 0.5um per dL",
			[]string{"serialno"},
			nil,
		),
		pm50CountDesc: prometheus.NewDesc(
			"airgradient_pm50_count",
			"Particle count 5um per dL",
			[]string{"serialno"},
			nil,
		),
		pm10CountDesc: prometheus.NewDesc(
			"airgradient_pm10_count",
			"Particle count 10um per dL",
			[]string{"serialno"},
			nil,
		),
		atmpDesc: prometheus.NewDesc(
			"airgradient_atmp",
			"Temperature in Degrees Celsius",
			[]string{"serialno"},
			nil,
		),
		atmpCompensatedDesc: prometheus.NewDesc(
			"airgradient_atmp_compensated",
			"Temperature in Degrees Celsius with correction applied",
			[]string{"serialno"},
			nil,
		),
		rhumDesc: prometheus.NewDesc(
			"airgradient_rhum",
			"Relative Humidity",
			[]string{"serialno"},
			nil,
		),
		rhumCompensatedDesc: prometheus.NewDesc(
			"airgradient_rhum_compensated",
			"Relative Humidity with correction applied",
			[]string{"serialno"},
			nil,
		),
		tvocIndexDesc: prometheus.NewDesc(
			"airgradient_tvoc_index",
			"Senisiron VOC Index",
			[]string{"serialno"},
			nil,
		),
		tvocRawDesc: prometheus.NewDesc(
			"airgradient_tvoc_raw",
			"VOC raw value",
			[]string{"serialno"},
			nil,
		),
		noxIndexDesc: prometheus.NewDesc(
			"airgradient_nox_index",
			"Senisirion NOx Index",
			[]string{"serialno"},
			nil,
		),
		noxRawDesc: prometheus.NewDesc(
			"airgradient_nox_raw",
			"NOx raw value",
			[]string{"serialno"},
			nil,
		),
		bootDesc: prometheus.NewDesc(
			"airgradient_boot_total",
			"The total uptime of the device in minutes",
			[]string{"serialno"},
			nil,
		),
	}, nil
}

type AirgradientCollector struct {
	ctx      context.Context
	client   *http.Client
	endpoint *url.URL

	deviceInfoDesc      *prometheus.Desc
	wifiDesc            *prometheus.Desc
	pm01Desc            *prometheus.Desc
	pm01StandardDesc    *prometheus.Desc
	pm01CountDesc       *prometheus.Desc
	pm02Desc            *prometheus.Desc
	pm02StandardDesc    *prometheus.Desc
	pm02CountDesc       *prometheus.Desc
	pm10Desc            *prometheus.Desc
	pm10StandardDesc    *prometheus.Desc
	pm02CompensatedDesc *prometheus.Desc
	rco2Desc            *prometheus.Desc
	pm003CountDesc      *prometheus.Desc
	pm005CountDesc      *prometheus.Desc
	pm50CountDesc       *prometheus.Desc
	pm10CountDesc       *prometheus.Desc
	atmpDesc            *prometheus.Desc
	atmpCompensatedDesc *prometheus.Desc
	rhumDesc            *prometheus.Desc
	rhumCompensatedDesc *prometheus.Desc
	tvocIndexDesc       *prometheus.Desc
	tvocRawDesc         *prometheus.Desc
	noxIndexDesc        *prometheus.Desc
	noxRawDesc          *prometheus.Desc
	bootDesc            *prometheus.Desc
}

func (c *AirgradientCollector) Describe(ch chan<- *prometheus.Desc) {
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
}

func (c *AirgradientCollector) Collect(ch chan<- prometheus.Metric) {
	m, err := c.GetMeasures(c.ctx)
	if err != nil {
		ilog.FromContext(c.ctx).Error("Failed to get measures.", zap.Error(err))
		return
	}

	ch <- prometheus.MustNewConstMetric(c.deviceInfoDesc, prometheus.GaugeValue, 1, m.SerialNo, m.Firmware, m.Model, m.LEDMode)
	ch <- prometheus.MustNewConstMetric(c.wifiDesc, prometheus.GaugeValue, float64(m.WiFi), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm01Desc, prometheus.GaugeValue, float64(m.PM01), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm01StandardDesc, prometheus.GaugeValue, float64(m.PM01Standard), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm01CountDesc, prometheus.GaugeValue, float64(m.PM01Count), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm02Desc, prometheus.GaugeValue, float64(m.PM02), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm02StandardDesc, prometheus.GaugeValue, float64(m.PM02Standard), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm02CountDesc, prometheus.GaugeValue, float64(m.PM02Count), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm10Desc, prometheus.GaugeValue, float64(m.PM10), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm10StandardDesc, prometheus.GaugeValue, float64(m.PM10Standard), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm02CompensatedDesc, prometheus.GaugeValue, float64(m.PM02Compensated), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.rco2Desc, prometheus.GaugeValue, float64(m.RCO2), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm003CountDesc, prometheus.GaugeValue, float64(m.PM003Count), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm005CountDesc, prometheus.GaugeValue, float64(m.PM005Count), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm50CountDesc, prometheus.GaugeValue, float64(m.PM50Count), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.pm10CountDesc, prometheus.GaugeValue, float64(m.PM10Count), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.atmpDesc, prometheus.GaugeValue, m.ATMP, m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.atmpCompensatedDesc, prometheus.GaugeValue, m.ATMPCompensated, m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.rhumDesc, prometheus.GaugeValue, float64(m.RHUM), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.rhumCompensatedDesc, prometheus.GaugeValue, float64(m.RHUMCompensated), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.tvocIndexDesc, prometheus.GaugeValue, float64(m.TVOCIndex), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.tvocRawDesc, prometheus.GaugeValue, float64(m.TVOCRaw), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.noxIndexDesc, prometheus.GaugeValue, float64(m.NOXIndex), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.noxRawDesc, prometheus.GaugeValue, float64(m.NOXRaw), m.SerialNo)
	ch <- prometheus.MustNewConstMetric(c.bootDesc, prometheus.CounterValue, float64(m.Boot), m.SerialNo)
}

// GetMeasures fetches the measures from the AirGradient local server API.
func (c *AirgradientCollector) GetMeasures(ctx context.Context) (*Measures, error) {
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

	var m Measures
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	ilog.FromContext(ctx).Debug("Got measures from airgradient.", zap.Any("measures", m))
	return &m, nil
}
