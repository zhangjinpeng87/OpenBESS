package balancing

import (
	"github.com/zhangjinpeng87/openbms/pkg/datamodel"
)

const (
	DefaultSOCShardCnt = 16
)

type PackSOC struct {
	// Accumulated data of all cells in the pack
	maxCapacity float64
	currentCapacity float64

	// cell id -> cell SOC
	cellSOC map[int]*datamodel.BatteryState
}

func NewPackSOC() *PackSOC {
	return &PackSOC{
		cellSOC: make(map[int]*datamodel.BatteryState),
	}
}

func (p *PackSOC) Update(state *datamodel.BatteryState) {
	p.cellSOC[state.Cell] = state
}

func (p *PackSOC) ReCalculate() float64 {
	maxCapacity, currentCapacity := 0.0, 0.0
	for _, cellSOC := range p.cellSOC {
		maxCapacity += cellSOC.MaxCapacity
		currentCapacity += cellSOC.CurrentCapacity
	}
	p.maxCapacity = maxCapacity
	p.currentCapacity = currentCapacity
}

type ContainerSOC struct {
	// Accumulated data of all packs in the container
	maxCapacity float64
	currentCapacity float64

	// pack id -> pack SOC
	packSOC map[int]*PackSOC
}

func NewContainerSOC() *ContainerSOC {
	return &ContainerSOC{
		packSOC: make(map[int]*PackSOC),
	}
}

func (c *ContainerSOC) Update(state *datamodel.BatteryState) {
	if _, ok := c.packSOC[state.Pack]; !ok {
		c.packSOC[state.Pack] = NewPackSOC()
	}
	c.packSOC[state.Pack].Update(state)
}

func (c *ContainerSOC) ReCalculate() float64 {
	maxCapacity, currentCapacity := 0.0, 0.0
	for _, packSOC := range c.packSOC {
		packSOC.ReCalculate()

		maxCapacity += packSOC.maxCapacity
		currentCapacity += packSOC.currentCapacity
	}
	c.maxCapacity = maxCapacity
	c.currentCapacity = currentCapacity
}

type StationSOC struct {
	// Accumulated data of all containers in the station
	maxCapacity float64
	currentCapacity float64

	// container id -> container SOC
	containerSOC map[int]*ContainerSOC
}

func NewStationSOC() *StationSOC {
	return &StationSOC{
		containerSOC: make(map[int]*ContainerSOC),
	}
}

func (s *StationSOC) Update(state *datamodel.BatteryState) {
	if _, ok := s.containerSOC[state.Container]; !ok {
		s.containerSOC[state.Container] = NewContainerSOC()
	}
	s.containerSOC[state.Container].Update(state)
}

func (s *StationSOC) ReCalculate() float64 {
	maxCapacity, currentCapacity := 0.0, 0.0
	for _, containerSOC := range s.containerSOC {
		containerSOC.ReCalculate()

		maxCapacity += containerSOC.maxCapacity
		currentCapacity += containerSOC.currentCapacity
	}
	s.maxCapacity = maxCapacity
	s.currentCapacity = currentCapacity
}

type SOCDataShard struct {
	// RWMutex to protect stationSOC from concurrent modification
	mu sync.RWMutex

	// Accumulated data of all stations
	maxCapacity float64
	currentCapacity float64

	// station id -> station SOC
	stationSOC map[int]*StationSOC
}

func NewSOCDataShard() *SOCDataShard {
	return &SOCDataShard{
		stationSOC: make(map[int]*StationSOC),
	}
}

func (s *SOCDataShard) Update(state *datamodel.BatteryState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.stationSOC[state.Station]; !ok {
		s.stationSOC[state.Station] = NewStationSOC()
	}
	s.stationSOC[state.Station].Update(state)
}

func (s *SOCDataShard) ReCalculate() {
	s.mu.Lock()
	defer s.mu.Unlock()

	maxCapacity, currentCapacity := 0.0, 0.0
	for _, stationSOC := range s.stationSOC {
		stationSOC.ReCalculate()

		maxCapacity += stationSOC.maxCapacity
		currentCapacity += stationSOC.currentCapacity
	}

	s.maxCapacity = maxCapacity
	s.currentCapacity = currentCapacity
}


type SOCData struct {
	shardCnt int
	shards []*SOCDataShard

	// Accumulated data of all shards
	maxCapacity float64
	currentCapacity float64
}

func NewSOCData(shardCnt int) *SOCData {
	shards := make([]*SOCDataShard, shardCnt)
	for i := 0; i < shardCnt; i++ {
		shards[i] = NewSOCDataShard()
	}
	return &SOCData{
		shardCnt: shardCnt,
		shards: shards,
	}
}

func (s *SOCData) Update(state *datamodel.BatteryState) {
	shard := s.shards[state.Station % s.shardCnt]
	shard.Update(state)
}

func (s *SOCData) ReCalculate() {
	for _, shard := range s.shards {
		shard.ReCalculate()
	}
}
