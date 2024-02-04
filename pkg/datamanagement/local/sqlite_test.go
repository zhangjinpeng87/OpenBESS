package localstore

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestSqliteStore_OpenClose(t *testing.T) {
	// Create a new SqliteStore with a mocked SQL database
	cfg := &config.LocalStoreConfig{Path: ":memory:"}
	store := NewSqliteStore(cfg)

	// Mock the SQL database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	store.db = sqlx.NewDb(mockDB, "sqlite3")

	// Expectations for the Open method
	mock.ExpectConnect()
	mock.ExpectExec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS battery_state.*")).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Call the Open method
	if err := store.Open(); err != nil {
		t.Errorf("Error opening database: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Expectations for the Close method
	mock.ExpectClose()

	// Call the Close method
	if err := store.Close(); err != nil {
		t.Errorf("Error closing database: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSqliteStore_Upsert(t *testing.T) {
	// Create a new SqliteStore with a mocked SQL database
	cfg := &config.LocalStoreConfig{Path: ":memory:"}
	store := NewSqliteStore(cfg)

	// Mock the SQL database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	store.db = sqlx.NewDb(mockDB, "sqlite3")

	// Sample BatteryState
	state := &datamodel.BatteryState{
		Station:        1,
		Container:      2,
		Pack:           3,
		Cell:           4,
		Voltage:        12.3,
		Current:        4.5,
		SOC:            78.9,
		Temperature:    25.5,
		State:          1,
		Timestamp:      time.Now().Unix(),
	}

	// Expectations for the Upsert method
	mock.ExpectExec(fmt.Sprintf("INSERT INTO battery_state.*")).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the Upsert method
	if err := store.Upsert(state); err != nil {
		t.Errorf("Error upserting battery state: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSqliteStore_GetLatest(t *testing.T) {
	// Create a new SqliteStore with a mocked SQL database
	cfg := &config.LocalStoreConfig{Path: ":memory:"}
	store := NewSqliteStore(cfg)

	// Mock the SQL database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	store.db = sqlx.NewDb(mockDB, "sqlite3")

	// Sample BatteryState
	state := &BatteryState{
		Station:     1,
		Container:   2,
		Pack:        3,
		Cell:        4,
		Voltage:     12.3,
		Current:     4.5,
		SOC:         78.9,
		Temperature: 25.5,
		State:       1,
		Timestamp:   time.Now().Unix(),
	}

	// Expectations for the GetLatest method
	rows := sqlmock.NewRows([]string{"station", "container", "pack", "cell", "voltage", "current", "soc", "temperature", "state", "timestamp"}).
		AddRow(state.Station, state.Container, state.Pack, state.Cell, state.Voltage, state.Current, state.SOC, state.Temperature, state.State, state.Timestamp)

	mock.ExpectQuery(fmt.Sprintf("SELECT \\* FROM battery_state.*")).
		WithArgs(state.Station, state.Container, state.Pack, state.Cell).
		WillReturnRows(rows)

	// Call the GetLatest method
	result, err := store.GetLatest(state.Station, state.Container, state.Pack, state.Cell)
	if err != nil {
		t.Errorf("Error getting latest battery state: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Perform assertions
	if result == nil {
		t.Error("Expected a BatteryState, got nil")
	} else {
		// Perform additional assertions if needed
	}
}

func TestSqliteStore_GenerateSnapshotFile_CSV(t *testing.T) {
	// Create a new SqliteStore with a mocked SQL database
	cfg := &config.LocalStoreConfig{Path: ":memory:"}
	store := NewSqliteStore(cfg)

	// Mock the SQL database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	store.db = sqlx.NewDb(mockDB, "sqlite3")

	// Sample BatteryStates
	states := []BatteryState{
		{Station: 1, Container: 1, Pack: 1, Cell: 1, Voltage: 12.3, Current: 4.5, SOC: 78.9, Temperature: 25.5, State: 1, Timestamp: time.Now().Unix()},
		{Station: 1, Container: 1, Pack: 1, Cell: 2, Voltage: 11.8, Current: 3.7, SOC: 82.1, Temperature: 26.3, State: 1, Timestamp: time.Now().Unix()},
		// Add more states as needed
	}

	// Expectations for the GenerateSnapshotFile method with CSV format
	mock.ExpectQuery("SELECT \\* FROM battery_state.*").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"station", "container", "pack", "cell", "voltage", "current", "soc", "temperature", "state", "timestamp"}).
			AddRow(states[0].Station, states[0].Container, states[0].Pack, states[0].Cell, states[0].Voltage, states[0].Current, states[0].SOC, states[0].Temperature, states[0].State, states[0].Timestamp).
			AddRow(states[1].Station, states[1].Container, states[1].Pack, states[1].Cell, states[1].Voltage, states[1].Current, states[1].SOC, states[1].Temperature, states[1].State, states[1].Timestamp))

	// Call the GenerateSnapshotFile method with CSV format
	file, _, err := store.GenerateSnapshotFile(1, "csv")
	if err != nil {
		t.Errorf("Error generating CSV snapshot file: %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

	// Perform assertions
	fileInfo, err := os.Stat(file)
	if err != nil {
		t.Errorf("Error getting file info: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Error("Expected non-empty CSV file, got empty file")
	}

	// Clean up (remove the generated file)
	if err := os.Remove(file); err != nil {
		t.Errorf("Error cleaning up: %v", err)
	}
}
