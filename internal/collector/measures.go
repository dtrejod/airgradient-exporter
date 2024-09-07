package collector

const (
	measuresPath = "/measures/current"
)

type measures struct {
	SerialNo        string  `json:"serialno"`
	Wifi            int     `json:"wifi"`
	PM01            int     `json:"pm01"`
	PM02            int     `json:"pm02"`
	PM10            int     `json:"pm10"`
	PM02Compensated int     `json:"pm02Compensated"`
	RCO2            int     `json:"rco2"`
	PM003Count      int     `json:"pm003Count"`
	ATMP            float64 `json:"atmp"`
	ATMPCompensated float64 `json:"atmpCompensated"`
	RHUM            int     `json:"rhum"`
	RHUMCompensated int     `json:"rhumCompensated"`
	TVOCIndex       int     `json:"tvocIndex"`
	TVOCRaw         int     `json:"tvocRaw"`
	NOXIndex        int     `json:"noxIndex"`
	NOXRaw          int     `json:"noxRaw"`
	Boot            int     `json:"boot"`
	BootCount       int     `json:"bootCount"`
	LEDMode         string  `json:"ledMode"`
	Firmware        string  `json:"firmware"`
	Model           string  `json:"model"`
}
