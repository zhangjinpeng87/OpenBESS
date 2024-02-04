package data_model

import (
	"testing"
)

func TestPackData(t *testing.T) {
	// Create a new PackData and update it with some test data
	packData := NewPackData()
	state := &datamodel.BatteryState{
		Cell:            1,
		MaxCapacity:     100.0,
		CurrentCapacity: 50.0,
		Pack:            1,
		Container:       1,
		Station:         1,
	}
	packData.Update(state)

	// Recalculate the PackData
	packData.ReCalculate()

	// Perform assertions
	if packData.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", packData.maxCapacity)
	}
	if packData.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", packData.currentCapacity)
	}
}

func TestContainerData(t *testing.T) {
	containerData := NewContainerData()

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update ContainerData with the test BatteryState
	containerData.Update(state)

	// Recalculate ContainerData
	containerData.ReCalculate()

	// Perform assertions
	if containerData.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", containerData.maxCapacity)
	}
	if containerData.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", containerData.currentCapacity)
	}
}

func TestStationData(t *testing.T) {
	stationData := NewStationData()

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update StationData with the test BatteryState
	stationData.Update(state)

	// Recalculate StationData
	stationData.ReCalculate()

	// Perform assertions
	if stationData.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", stationData.maxCapacity)
	}
	if stationData.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", stationData.currentCapacity)
	}
}

func TestDataDataShard(t *testing.T) {
	DataDataShard := NewDataDataShard()

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update DataDataShard with the test BatteryState
	DataDataShard.Update(state)

	// Recalculate DataDataShard
	DataDataShard.ReCalculate()

	// Perform assertions
	if DataShard.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", DataShard.maxCapacity)
	}
	if DataShard.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", DataShard.currentCapacity)
	}
}

func TestBatteriesData(t *testing.T) {
	shardCnt := DefaultShardCnt
	BatteriesData := NewBatteriesData(shardCnt)

	// Create a test BatteryState
	state := &datamodel.BatteryState{
		Cell:           1,
		MaxCapacity:    100.0,
		CurrentCapacity: 50.0,
		Pack:           1,
		Container:      1,
		Station:        1,
	}

	// Update Data with the test BatteryState
	BatteriesData.Update(state)

	// Recalculate Data
	BatteriesData.ReCalculate()

	// Perform assertions
	if BatteriesData.maxCapacity != 100.0 {
		t.Errorf("Expected maxCapacity to be 100.0, got %f", BatteriesData.maxCapacity)
	}
	if BatteriesData.currentCapacity != 50.0 {
		t.Errorf("Expected currentCapacity to be 50.0, got %f", BatteriesData.currentCapacity)
	}
}
