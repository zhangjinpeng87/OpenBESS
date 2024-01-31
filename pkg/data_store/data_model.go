package datamodel

type State int

const (
	Idle State = iota + 1
	Charging
	Discharging
	FastCharging
	FastDischarging
)

func (s State) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Charging:
		return "Charging"
	case Discharging:
		return "Discharging"
	case FastCharging:
		return "FastCharging"
	case FastDischarging:
		return "FastDischarging"
	default:
		return "Unknown"
	}
}

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
	// SOC is the estimated state of charge of the battery.
	SOC float64 `json:"soc"`
	// SOH is the estimated state of health of the battery.
	SOH float64 `json:"soh"`
	// MaxCapacity is the maximum capacity of the battery in ampere hours.
	MaxCapacity float64 `json:"max_capacity"`
	// Temperature is the battery temperature in degrees Celsius.
	Temperature float64 `json:"temperature"`
	// Timestamp is the timestamp of the battery state.
	Timestamp int64 `json:"timestamp"`
	// State is the state of the battery.
	State State `json:"state"`
}
