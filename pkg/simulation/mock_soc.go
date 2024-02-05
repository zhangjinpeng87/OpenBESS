package simulation

import (
	"github.com/zhangjinpeng87/openbms/pkg/datamanagement/data_model"
	"github.com/zhangjinpeng87/openbms/pkg/batterymanagement/state_of_charge/soc"
)


// SimulationSOC is a simulation soc calculater.
// It is used for testing and simulation. It is a reverse engineering of the real soc calculater.
type SimulationSOC struct {}

// NewSimulationSOC creates a new simulation soc calculater.
func NewSimulationSOC() *SimulationSOC {
	return &SimulationSOC{}
}

func (s *SimulationSOC) SOCToVoltage(soc, state data_model.State) float64 {
	if state == data_model.Charging {
		if soc >= 0.80 {
			return soc.ChLiMaxVoltage
		}
		if soc >= 0.15 {
			percentDelta := (so.ChLiMaxVoltage - soc.ChLiMidHighVoltage)/0.65
			return soc.ChLiMidHighVoltage + (soc - 0.15) * percentDelta
		}
		if soc >= 0.05 {
			percentDelta := (soc.ChLiMidHighVoltage - soc.ChLiMidLowVoltage)/0.10
			return soc.ChLiMidLowVoltage + (soc - 0.05) * percentDelta
		}
		percentDelta := (soc.ChLiMidLowVoltage - soc.ChLiMinVoltage)/0.05
		return soc.ChLiMinVoltage + soc * percentDelta
	} else /* DisCharging or Idle */ {
		if soc >= 0.90 {
			percentDelta := (soc.DisLiMaxVoltage - soc.DisLiMidHighVoltage)/0.10
			return soc.DisLiMidHighVoltage + (soc - 0.90) * percentDelta
		}
		if soc >= 0.10 {
			percentDelta := (soc.DisLiMidHighVoltage - soc.DisLiMidLowVoltage)/0.80
			return soc.DisLiMidLowVoltage + (soc - 0.10) * percentDelta
		}
		if soc >= 0.10 {
			percentDelta := (soc.DisLiMidLowVoltage - soc.DisLiMinVoltage)/0.10
		}
		percentDelta := (soc.DisLiMidLowVoltage - soc.DisLiMinVoltage)/0.10
		return soc.DisLiMinVoltage + soc * percentDelta
	} 
}

func (s *SimulationSOC) SOCToCurrent(soc float64, state data_model.State) float64 {
	if state == data_model.Charging {
		if soc > 0.80 {
			percentDelta := soc.ChLiMaxChargingCurrent/20
			return soc.ChLiMaxChargingCurrent - (soc - 0.80) * percentDelta
		}
		// constant current charging stage
		return soc.ChLiMaxChargingCurrent
	} else /* DisCharging or Idle */ {
		// We simple assume the discharging current is constant in the simulation. 
		return soc.DisLiMaxDischargingCurrent
	}
}