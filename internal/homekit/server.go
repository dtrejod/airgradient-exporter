package homekit

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"path/filepath"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	hclog "github.com/brutella/hc/log"
	"github.com/dtrejod/airgradient-exporter/internal/collector"
	"github.com/dtrejod/airgradient-exporter/internal/ilog"
	"github.com/dtrejod/airgradient-exporter/version"
	"go.uber.org/zap"
)

const (
	homekitPairPin = "18458232"

	co2AbnormalThreshold = 800 // ppm

	// Thresholds from EPA. See table "2024 AQI for Fine Particle Pollution" from reference.
	// REF: https://www.epa.gov/system/files/documents/2024-02/pm-naaqs-air-quality-index-fact-sheet.pdf
	p02ExcellentThreshold = float64(5)     // µg/m³
	p02GoodThreshold      = float64(9.0)   // µg/m³
	p02FairThreshold      = float64(35.4)  // µg/m³
	p02InferiorThreshold  = float64(55.4)  // µg/m³
	p02PoorThreshold      = float64(125.4) // µg/m³
)

func NewServer(ctx context.Context, listenAddr, dataDir string, ag *collector.AirgradientCollector) (*hap.Server, error) {
	// Disable logging since it does not use zap structured logging.
	hclog.Info.Disable()
	hclog.Debug.Disable()

	bridge := accessory.NewBridge(accessory.Info{
		Name:     "AirGradient Exporter",
		Firmware: version.Version(),
	})

	// Create the AirGradient accessory.
	a, err := newAirgradientAccessory(ctx, ag)
	if err != nil {
		return nil, fmt.Errorf("could not create thermometer accessory: %w", err)
	}

	// Store the data in the "./data/db" directory.
	fs := hap.NewFsStore(filepath.Join(dataDir, "db"))

	// Create the hap server.
	server, err := hap.NewServer(fs, bridge.A, a)
	if err != nil {
		return nil, fmt.Errorf("could not start hap (homekit) server: %w", err)
	}
	// Create a random pin
	server.Pin = homekitPairPin
	server.Addr = listenAddr // Set the listen address if provided.

	ilog.FromContext(ctx).Info("Homekit server created successfully", zap.String("pin", homekitPairPin))
	return server, nil
}

func newAirgradientAccessory(ctx context.Context, ag *collector.AirgradientCollector) (*accessory.A, error) {
	measures, err := ag.GetMeasures(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get measures from airgradient collector: %w", err)
	}

	sensor := accessory.New(accessory.Info{
		Name:         "AirGradient ONE",
		Manufacturer: "AirGradient",
		SerialNumber: measures.SerialNo,
		Model:        measures.Model,
		Firmware:     measures.Firmware,
	}, accessory.TypeSensor)
	hasher := fnv.New64()
	hasher.Write([]byte(measures.Model))
	hasher.Write([]byte(measures.SerialNo))
	// Set the unique ID for the accessory.
	sensor.Id = hasher.Sum64()

	// Creator sensors services compatible with HomeKit.
	// Ref: https://developer.apple.com/documentation/homekit/accessory-service-types#Temperature-and-Humidity
	temp := service.NewTemperatureSensor()
	humidity := service.NewHumiditySensor()
	co2 := service.NewCarbonDioxideSensor()
	quality := service.NewAirQualitySensor()

	temp.CurrentTemperature.ValueRequestFunc = func(req *http.Request) (any, int) {
		ilog.FromContext(ctx).Debug("Requesting temperature value from AirGradient collector")
		measures, err := ag.GetMeasures(ctx)
		if err != nil {
			ilog.FromContext(ctx).Error("Failed to get temperature from AirGradient collector", zap.Error(err))
			return nil, hap.JsonStatusServiceCommunicationFailure
		}
		return measures.ATMP, hap.JsonStatusSuccess
	}
	humidity.CurrentRelativeHumidity.ValueRequestFunc = func(req *http.Request) (any, int) {
		ilog.FromContext(ctx).Debug("Getting humidity from AirGradient collector")
		measures, err := ag.GetMeasures(ctx)
		if err != nil {
			ilog.FromContext(ctx).Error("Failed to get humidity from AirGradient collector", zap.Error(err))
			return nil, hap.JsonStatusServiceCommunicationFailure
		}
		return measures.RHUM, hap.JsonStatusSuccess
	}
	co2.CarbonDioxideDetected.ValueRequestFunc = func(req *http.Request) (any, int) {
		ilog.FromContext(ctx).Debug("Requesting CO2 level from AirGradient collector")
		measures, err := ag.GetMeasures(ctx)
		if err != nil {
			ilog.FromContext(ctx).Error("Failed to get CO2 from AirGradient collector", zap.Error(err))
			return nil, hap.JsonStatusServiceCommunicationFailure
		}
		if measures.RCO2 > co2AbnormalThreshold {
			return characteristic.CarbonDioxideDetectedCO2LevelsAbnormal, hap.JsonStatusSuccess
		}
		return characteristic.CarbonDioxideDetectedCO2LevelsNormal, hap.JsonStatusSuccess
	}
	quality.AirQuality.ValueRequestFunc = func(req *http.Request) (any, int) {
		ilog.FromContext(ctx).Debug("Getting air quality from AirGradient collector")
		measures, err := ag.GetMeasures(ctx)
		if err != nil {
			ilog.FromContext(ctx).Error("Failed to get air quality from AirGradient collector", zap.Error(err))
			return nil, hap.JsonStatusServiceCommunicationFailure
		}
		if measures.PM02 < 0 {
			return characteristic.AirQualityUnknown, hap.JsonStatusSuccess
		} else if measures.PM02 < p02ExcellentThreshold {
			return characteristic.AirQualityExcellent, hap.JsonStatusSuccess
		} else if measures.PM02 < p02GoodThreshold {
			return characteristic.AirQualityGood, hap.JsonStatusSuccess
		} else if measures.PM02 < p02FairThreshold {
			return characteristic.AirQualityFair, hap.JsonStatusSuccess
		} else if measures.PM02 < p02InferiorThreshold {
			return characteristic.AirQualityInferior, hap.JsonStatusSuccess
		} else if measures.PM02 < p02PoorThreshold {
			return characteristic.AirQualityPoor, hap.JsonStatusSuccess
		}
		return characteristic.AirQualityUnknown, hap.JsonStatusSuccess
	}

	// Add the services to the accessory.
	sensor.AddS(temp.S)
	sensor.AddS(humidity.S)
	sensor.AddS(co2.S)
	sensor.AddS(quality.S)

	return sensor, nil
}
