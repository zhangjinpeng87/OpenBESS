package soc

import (
	"github.com/zhangjinpeng87/openbms/pkg/utils"
)

// State Of Charge (SOC) is the available capacity of a battery expressed as a percentage of the rated capacity.
// Use Kalman Filter to filter the noise of voltage and current.
//  1.1 Discharge Voltage Curve
//  V |
//    |  .
//    |    .
//    |      .
//    |       .   .     .     .    .    .    .   . .
//    |                                              .
//    |    									          .
//    |   									            .
//    |___________________________________________________________
//    100%  90%  80%  70%  60%  50%  40%  30%  20%  10%  0% (Capacity)
//    |------|---------------------------------------|-----|
//   4.2V   3.7V                                    3.0V  2.5V

//  2.1 Charge Voltage Curve
//  V |
//    |                   .  . 	 .     .    .    .    .   .   .
//    |      .
//    |    .
//    |   .
//    |  .
//    | .
//    |.
//    |___________________________________________________________
//    0    10    20    30    40    50    60    70 (Charge Time in Minutes)
//
//  2.2 Charge Current Curve
//  I |
//    |  .   .   .   .   .
//    |
//    |                    .
//    |                      .
//    |                        .
//    |                          .
//    |                            .
//    |                               .  .  .  .
//    |___________________________________________________________
//    0    10    20    30    40    50    60    70 (Charge Time in Minutes)
//    |------------------|----------------------|
//   Constant Current Charge   Saturation Charge

const (
	// Discharge parameters of Li-ion battery
	// DisLiMaxVoltage is the maximum voltage of a battery cell when discharging.
	DisLiMaxVoltage = 4.25
	// DisLiMidHighVoltage is the mid high voltage of a battery cell when discharging.
	DisLiMidHighVoltage = 3.7
	// DisLiMidLowVoltage is the mid low voltage of a battery cell when discharging.
	DisLiMidLowVoltage = 3.0
	// DisLiMinVoltage is the minimum voltage of a battery cell when discharging.
	DisLiMinVoltage = 2.5

	// Charge parameters of Li-ion battery
	// ChLiMaxVoltage is the maximum voltage of a battery cell when charging.
	ChLiMaxVoltage = 4.2
	// ChLiMidHighVoltage is the mid high voltage of a battery cell when charging.
	ChLiMidHighVoltage = 3.8
	// ChLiMidLowVoltage is the mid low voltage of a battery cell when charging.
	ChLiMidLowVoltage = 3.5
	// ChLiMinVoltage is the minimum voltage of a battery cell when charging.
	ChLiMinVoltage = 2.5
	// LiMaxChargingCurrent is the maximum charging current of a battery cell.
	LiMaxChargingCurrent = 1.0

	VoltageDeviation = 0.1
)

type DischargeCalculater struct {
	DisMaxVoltage     float64
	DisMidHighVoltage float64
	DisMidLowVoltage  float64
	DisMinVoltage     float64

	KalmanFilter *utils.KalmanFilter
}

// NewSocCalculater creates a new soc calculater.
func DischargeCalculater(MaxVoltage, MidHighVoltage, MidLowVoltage, MinVoltage float64) *DischargeCalculater {
	return &DischargeSocCalculater{
		DisMaxVoltage:     MaxVoltage,
		DisMidHighVoltage: MidHighVoltage,
		DisMidLowVoltage:  MidLowVoltage,
		DisMinVoltage:     MinVoltage,
	}
}

func (s *DischargeCalculater) InitKalmanFilter(initialValue, initialEstimateError, processNoise, measurementNoise float64) {
	s.KalmanFilter = utils.NewKalmanFilter(initialValue, initialEstimateError, processNoise, measurementNoise)
}

// SOC calculates the soc of a battery cell when discharging.
func (s *DischargeCalculater) SOC(voltage float64) float64 {
	if s.KalmanFilter != nil {
		voltage = s.KalmanFilter.Update(voltage)
	}

	if voltage >= s.DisMaxVoltage {
		return 100
	}
	if voltage >= s.DisMidHighVoltage {
		return 100 - (s.DisMaxVoltage-voltage)/(s.DisMaxVoltage-s.DisMidHighVoltage)*10
	}
	if voltage >= s.DisMidLowVoltage {
		return 90 - (s.DisMidHighVoltage-voltage)/(s.DisMidHighVoltage-s.DisMidLowVoltage)*80
	}
	if voltage >= s.DisMinVoltage {
		return 10 - (s.DisMidLowVoltage-voltage)/(s.DisMidLowVoltage-s.DisMinVoltage)*10
	}
	return 0
}

type ChargeCalculater struct {
	ChMaxVoltage         float64
	ChMidHighVoltage     float64
	ChMidLowVoltage      float64
	ChMinVoltage         float64
	ChMaxChargingCurrent float64

	// Voltage Kalman Filter
	VKalmanFilter *utils.KalmanFilter
	// Current Kalman Filter
	CKalmanFilter *utils.KalmanFilter
}

// NewChargeCalculater creates a new charge calculater.
func NewChargeCalculater(MaxVoltage, MidHighVoltage, MidLowVoltage, MinVoltage, MaxChargingCurrent float64) *ChargeCalculater {
	return &ChargeSocCalculater{
		ChMaxVoltage:         MaxVoltage,
		ChMidHighVoltage:     MidHighVoltage,
		ChMidLowVoltage:      MidLowVoltage,
		ChMinVoltage:         MinVoltage,
		ChMaxChargingCurrent: MaxChargingCurrent,
	}
}

func (s *ChargeCalculater) InitKalmanFilter(initialValue, initialEstimateError, processNoise, measurementNoise float64) {
	s.VKalmanFilter = utils.NewKalmanFilter(initialValue, initialEstimateError, processNoise, measurementNoise)
	s.CKalmanFilter = utils.NewKalmanFilter(initialValue, initialEstimateError, processNoise, measurementNoise)
}

// SOC calculates the soc of a battery cell when charging.
func (s *ChargeCalculater) SOC(voltage, current float64) float64 {
	if s.VKalmanFilter != nil {
		voltage = s.VKalmanFilter.Update(voltage)
	}
	if s.CKalmanFilter != nil {
		current = s.CKalmanFilter.Update(current)
	}

	if current >= s.MaxChargingCurrent {
		// Constant Current Charge Stage
		if voltage < s.ChMinVoltage {
			return 0
		} else if voltage >= s.ChMinVoltage && voltage < s.ChMidLowVoltage {
			return (voltage - s.ChMinVoltage) / (s.ChMidLowVoltage - s.ChMinVoltage) * 5
		} else if voltage >= s.ChMidLowVoltage && voltage < s.ChMidHighVoltage {
			return 5 + (voltage-s.ChMidLowVoltage)/(s.ChMidHighVoltage-s.ChMidLowVoltage)*10
		} else /* voltage >= s.ChMidHighVoltage */ {
			return 80 - (s.ChMaxVoltage-voltage)/(s.ChMaxVoltage-s.ChMidHighVoltage)*65
		}
	} else {
		// Saturation Charge Stage
		return 80 + (s.MaxChargingCurrent-current)/s.MaxChargingCurrent*20
	}
}
