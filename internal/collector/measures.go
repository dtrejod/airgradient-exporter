package collector

const (
	measuresPath = "/measures/current"
)

type measures struct {
	PM01            float64 `json:"pm01"`
	PM02            float64 `json:"pm02"`
	PM10            float64 `json:"pm10"`
	PM01Standard    float64 `json:"pm01Standard"`
	PM02Standard    float64 `json:"pm02Standard"`
	PM10Standard    float64 `json:"pm10Standard"`
	PM003Count      float64 `json:"pm003Count"`
	PM005Count      float64 `json:"pm005Count"`
	PM01Count       float64 `json:"pm01Count"`
	PM02Count       float64 `json:"pm02Count"`
	PM50Count       float64 `json:"pm50Count"`
	PM10Count       float64 `json:"pm10Count"`
	ATMP            float64 `json:"atmp"`
	ATMPCompensated float64 `json:"atmpCompensated"`
	RHUM            float64 `json:"rhum"`
	RHUMCompensated float64 `json:"rhumCompensated"`
	PM02Compensated float64 `json:"pm02Compensated"`
	RCO2            float64 `json:"rco2"`
	TVOCIndex       float64 `json:"tvocIndex"`
	TVOCRaw         float64 `json:"tvocRaw"`
	NOXIndex        int     `json:"noxIndex"`
	NOXRaw          float64 `json:"noxRaw"`
	Boot            int     `json:"boot"`
	// Deprecated: BootCount is deprecated in favor of Boot
	BootCount int    `json:"bootCount"`
	WiFi      int    `json:"wifi"`
	LEDMode   string `json:"ledMode"`
	SerialNo  string `json:"serialno"`
	Firmware  string `json:"firmware"`
	Model     string `json:"model"`
}
