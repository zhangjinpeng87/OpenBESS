package simulation


import (
	"context"
	"time"
	"golang.org/x/sync/errgroup"

	"github.com/zhangjinpeng87/openbms/pkg/datamanagement/data_model"
)

// Sensor is a sensor.
type MockSensor struct {
	batteryState data_model.BatteryState

	// Report Channel
	reportCh chan<- data_model.BatteryState
}

func NewMockSensor(station, container, pack, cell int, reportCh chan<- data_model.BatteryState) *MockSensor {
	var initState data_model.BatteryState
	initState.Station = station
	initState.Container = container
	initState.Pack = pack
	initState.Cell = cell

	// Randomly initialize the battery state.
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	r := rand.Intn(3)
	initState.State = if r == 0 {
		data_model.Idle
	} if r == 1 {
		data_model.Charging
	} else {
		data_model.Discharging
	}
	initState.SOC = rand.Float64(0.1, 1)

	return &MocSensor{
		batteryState: initState,
		reportCh: reportCh,
	}
}

func (s *MockSensor) Start(ctx context, reportInterval time.Duration, reportCh ) error {
	errg, c := errgroup.WithContext(ctx)
	errg.Go(func() error {
		timer := time.NewTimer(reportInterval)
		for {
			select {
			case <-c.Done():
				return nil
			case <-timer.C:
				s.Report()
				timer.Reset(reportInterval)
			}
		}
}

func (s *MockSensor) Report() {
	// Report the sensor data to the data management system.
	batteryState := s.batteryState
	batteryState.Timestamp = time.Now().Unix()
	
	// Simulate the voltage and current based on the SOC and state.
	simSoc := soc.NewSimulationSOC()
	voltage := simSoc.SOCToVoltage(batteryState.SOC, batteryState.State)
	current := simSoc.SOCToCurrent(batteryState.SOC, batteryState.State)
	batteryState.Voltage = voltage
	batteryState.Current = current

	// Report the battery state.
	s.reportCh <- batteryState
}
