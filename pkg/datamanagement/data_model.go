package data_model

import (
	"sync"

	"github.com/zhangjinpeng87/openbms/pkg/batterymanagement/state_of_charge/soc"
	"github.com/zhangjinpeng87/openbms/pkg/utils"
)

type State int

const (
	Idle State = iota + 1
	Charging
	Discharging
)

func (s State) String() string {
	switch s {
	case Idle:
		return "Idle"
	case Charging:
		return "Charging"
	case Discharging:
		return "Discharging"
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

const (
	DefaultDataShardCnt = 16
)

type PackData struct {
	// Accumulated data of all cells in the pack
	maxCapacity     float64
	currentCapacity float64

	// cell id -> cell Data
	cellData map[int]*datamodel.BatteryState

	// cell id -> cell voltage kalman filter
	cellKalmanV map[int]*utils.KalmanFilter
	// cell id -> cell current kalman filter
	cellKalmanC map[int]*utils.KalmanFilter
	// cell id -> cell temperature kalman filter
	cellKalmanT map[int]*utils.KalmanFilter

	// soc calculator
	dischargeSOCCalc *soc.DischargeSOCCalculator
	chargeSOCCalc    *soc.ChargeSOCCalculator
}

func NewPackData() *PackData {
	return &PackData{
		cellData:         make(map[int]*datamodel.BatteryState),
		dischargeSOCCalc: soc.NewDefaultDischargeSOCCalculator(),
		chargeSOCCalc:    soc.NewDefaultChargeSOCCalculator(),
	}
}

func (p *PackData) Update(state *datamodel.BatteryState) {
	// Smoonth the data collected by the sensors
	if _, ok := p.cellKalmanV[state.Cell]; !ok {
		p.cellKalmanV[state.Cell] = utils.NewKalmanFilter(state.Voltage, 1, 0.01, 0.01)
	}
	state.Voltage = p.cellKalmanV[state.Cell].Update(state.Voltage)
	if _, ok := p.cellKalmanC[state.Cell]; !ok {
		p.cellKalmanC[state.Cell] = utils.NewKalmanFilter(state.Current, 1, 0.01, 0.01)
	}
	state.Current = p.cellKalmanC[state.Cell].Update(state.Current)
	if _, ok := p.cellKalmanT[state.Cell]; !ok {
		p.cellKalmanT[state.Cell] = utils.NewKalmanFilter(state.Temperature, 1, 0.01, 0.01)
	}
	state.Temperature = p.cellKalmanT[state.Cell].Update(state.Temperature)

	p.cellData[state.Cell] = state
}

func (p *PackData) ReCalculate() float64 {
	maxCapacity, currentCapacity := 0.0, 0.0
	for _, cellData := range p.cellData {
		maxCapacity += cellData.MaxCapacity

		var socCalc *soc.SOCalculator
		if cellData.State == datamodel.Charging {
			socCalc = p.chargeSOCCalc
		} else {
			socCalc = p.dischargeSOCCalc
		}
		cellData.SOC = socCalc.Calculate(cellData.Voltage, cellData.Current)
		currentCapacity += cellData.SOC * cellData.MaxCapacity
	}
	p.maxCapacity = maxCapacity
	p.currentCapacity = currentCapacity
}

type ContainerData struct {
	// Accumulated data of all packs in the container
	maxCapacity     float64
	currentCapacity float64

	// pack id -> pack Data
	packData map[int]*PackData
}

func NewContainerData() *ContainerData {
	return &ContainerData{
		packData: make(map[int]*PackData),
	}
}

func (c *ContainerData) Update(state *datamodel.BatteryState) {
	if _, ok := c.packData[state.Pack]; !ok {
		c.packData[state.Pack] = NewPackData()
	}
	c.packData[state.Pack].Update(state)
}

func (c *ContainerData) ReCalculate() float64 {
	maxCapacity, currentCapacity := 0.0, 0.0
	for _, packData := range c.packData {
		packData.ReCalculate()

		maxCapacity += packData.maxCapacity
		currentCapacity += packData.currentCapacity
	}
	c.maxCapacity = maxCapacity
	c.currentCapacity = currentCapacity
}

type StationData struct {
	// Accumulated data of all containers in the station
	maxCapacity     float64
	currentCapacity float64

	// container id -> container Data
	containerData map[int]*ContainerData
}

func NewStationData() *StationData {
	return &StationData{
		containerData: make(map[int]*ContainerData),
	}
}

func (s *StationData) Update(state *datamodel.BatteryState) {
	if _, ok := s.containerData[state.Container]; !ok {
		s.containerData[state.Container] = NewContainerData()
	}
	s.containerData[state.Container].Update(state)
}

func (s *StationData) ReCalculate() float64 {
	maxCapacity, currentCapacity := 0.0, 0.0
	for _, containerData := range s.containerData {
		containerData.ReCalculate()

		maxCapacity += containerData.maxCapacity
		currentCapacity += containerData.currentCapacity
	}
	s.maxCapacity = maxCapacity
	s.currentCapacity = currentCapacity
}

type DataShard struct {
	// RWMutex to protect stationData from concurrent modification
	mu sync.RWMutex

	// Accumulated data of all stations
	maxCapacity     float64
	currentCapacity float64

	// station id -> station Data
	stationData map[int]*StationData
}

func NewDataShard() *DataShard {
	return &DataDataShard{
		stationData: make(map[int]*StationData),
	}
}

func (s *DataShard) Update(state *datamodel.BatteryState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.stationData[state.Station]; !ok {
		s.stationData[state.Station] = NewStationData()
	}
	s.stationData[state.Station].Update(state)
}

func (s *DataShard) ReCalculate() {
	s.mu.Lock()
	defer s.mu.Unlock()

	maxCapacity, currentCapacity := 0.0, 0.0
	for _, stationData := range s.stationData {
		stationData.ReCalculate()

		maxCapacity += stationData.maxCapacity
		currentCapacity += stationData.currentCapacity
	}

	s.maxCapacity = maxCapacity
	s.currentCapacity = currentCapacity
}

type BatteriesData struct {
	shardCnt int
	shards   []*DataShard

	// Accumulated data of all shards
	maxCapacity     float64
	currentCapacity float64
}

func NewBatteriesData(shardCnt int) *BatteriesData {
	shards := make([]*DataShard, shardCnt)
	for i := 0; i < shardCnt; i++ {
		shards[i] = NewDataShard()
	}
	return &BatteriesData{
		shardCnt: shardCnt,
		shards:   shards,
	}
}

func (s *BatteriesData) Update(state *datamodel.BatteryState) {
	shard := s.shards[state.Station%s.shardCnt]
	shard.Update(state)
}

func (s *BatteriesData) ReCalculate() {
	for _, shard := range s.shards {
		shard.ReCalculate()
	}
}
