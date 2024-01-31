package balancing

import (
	"testing"
	"github.com/zhangjinpeng87/openbms/pkg/datamodel"
)

func TestPackSOC(t *testing.T) {
	// Create a new PackSOC and update it with some test data
	packSOC := NewPackSOC()
	state := &datamodel.BatteryState{
		Cell:            1,
		MaxCapacity:     100.0,
		CurrentCapacity: 50.0,
		Pack:            1,
		Container:       1,
		Station:         1,
	}
	packSOC.Update(state)

	// Recalculate the PackSOC
	packSOC.ReCalculate()

	// Perform assertions
	if packSOC.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", packSOC.maxCapacity)
	}
	if packSOC.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", packSOC.currentCapacity)
	}
}

func TestContainerSOC(t *testing.T) {
	containerSOC := NewContainerSOC()

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update ContainerSOC with the test BatteryState
	containerSOC.Update(state)

	// Recalculate ContainerSOC
	containerSOC.ReCalculate()

	// Perform assertions
	if containerSOC.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", containerSOC.maxCapacity)
	}
	if containerSOC.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", containerSOC.currentCapacity)
	}
}

func TestStationSOC(t *testing.T) {
	stationSOC := NewStationSOC()

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update StationSOC with the test BatteryState
	stationSOC.Update(state)

	// Recalculate StationSOC
	stationSOC.ReCalculate()

	// Perform assertions
	if stationSOC.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", stationSOC.maxCapacity)
	}
	if stationSOC.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", stationSOC.currentCapacity)
	}
}

func TestSOCDataShard(t *testing.T) {
	socDataShard := NewSOCDataShard()

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update SOCDataShard with the test BatteryState
	socDataShard.Update(state)

	// Recalculate SOCDataShard
	socDataShard.ReCalculate()

	// Perform assertions
	if socDataShard.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", socDataShard.maxCapacity)
	}
	if socDataShard.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", socDataShard.currentCapacity)
	}
}

func TestSOCData(t *testing.T) {
	shardCnt := DefaultSOCShardCnt
	socData := NewSOCData(shardCnt)

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update SOCData with the test BatteryState
	socData.Update(state)

	// Recalculate SOCData
	socData.ReCalculate()

	// Perform assertions
	if socData.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", socData.maxCapacity)
	}
	if socData.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", socData.currentCapacity)
	}
}
