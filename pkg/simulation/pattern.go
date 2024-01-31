package simulation

import (

)

type MockCfg struct {
	// Behavior is the behavior of the sensor.
	Behavior Behavior `json:"behavior"`

	StartVoltage float64 `json:"start_voltage"`
	StartCurrent float64 `json:"start_current"`
	StartTemperature float64 `json:"start_temperature"`
	StartCapacity float64 `json:"capacity"`
	DegradationFactor float64 `json:"degradation_factor"`
}