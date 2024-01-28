package datamodel

// BatteryState is the state of a battery cell.
type BatteryState struct {
	// Station id of the battery station.
	Station int `json:"station"`
	// Container id of the battery container.
	Container int `json:"container"`
	// Pack id of the battery pack.
	Pack int `json:"pack"`
	// Cell is the cell id of the battery cell.
	Cell int `json:"cell"`

	// Voltage is the battery voltage in volts.
	Voltage float64 `json:"voltage"`
	// Current is the battery current in amps.
	Current float64 `json:"current"`
	// Temperature is the battery temperature in degrees Celsius.
	Temperature float64 `json:"temperature"`
	// Timestamp is the timestamp of the battery state.
	Timestamp int64 `json:"timestamp"`
}