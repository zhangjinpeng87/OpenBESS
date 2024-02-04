package localstore

import (
	"bufio"
	"database/sql"
	"fmt"
	"hash/crc32"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangjinpeng87/openbms/pkg/config"
	"github.com/zhangjinpeng87/openbms/pkg/datamodel"
)

type SqliteStore struct {
	cfg *config.LocalStoreConfig

	db *sqlx.DB
}

func NewSqliteStore(cfg *config.LocalStoreConfig) *SqliteStore {
	return &SqliteStore{cfg: cfg}
}

// Open opens the database.
func (s *SqliteStore) Open() error {
	// open database
	db, err := sqlx.Connect("sqlite3", s.cfg.Path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(4)

	// init schema
	if err := s.initSchema(); err != nil {
		return fmt.Errorf("failed to init schema: %w", err)
	}

	s.db = db
	return nil
}

// Close closes the database.
func (s *SqliteStore) Close() error {
	return s.db.Close()
}

// initSchema initializes the database schema if not exists.
func (s *SqliteStore) initSchema() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS battery_state (
			station INTEGER NOT NULL,
			container INTEGER NOT NULL,
			pack INTEGER NOT NULL,
			cell INTEGER NOT NULL,
			voltage REAL NOT NULL,
			current REAL NOT NULL,
			soc REAL NOT NULL,
			temperature REAL NOT NULL,
			state INTEGER NOT NULL,
			timestamp INTEGER NOT NULL,
			PRIMARY KEY (station, container, pack, cell)
		);
	`)
	return err
}

// Upsert upserts a battery state.
// TODO: use perpared statement to improve performance.
func (s *SqliteStore) Upsert(state *datamodel.BatteryState) error {
	_, err := s.db.Exec(`
		INSERT INTO battery_state(station, container, pack, cell, voltage, current, soc, temperature, state, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(station, container, pack, cell) DO UPDATE SET
			voltage = excluded.voltage,
			current = excluded.current,
			soc = excluded.soc,
			temperature = excluded.temperature,
			state = excluded.state
			timestamp = excluded.timestamp,
	`, state.Station, state.Container, state.Pack, state.Cell, state.Voltage, state.Current, state.soc, state.Temperature, state.State, state.Timestamp)
	return err
}

// GetLatest gets the latest battery state.
func (s *SqliteStore) GetLatest(station, container, pack, cell int) (*BatteryState, error) {
	var state BatteryState
	err := s.db.Get(&state, `
		SELECT * FROM battery_state
		WHERE station = ? AND container = ? AND pack = ? AND cell = ?
		ORDER BY timestamp DESC
		LIMIT 1
	`, station, container, pack, cell)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &state, err
}

// GenerateSnapshotFile generates snapshot file of the battery state for a specified station.
// The snapshot file will be upload to the cloud storage by certain frequency like per minute.
// Format can be "csv" or "parquet".
func (s *SqliteStore) GenerateSnapshotFile(station int, format string) (file string, checksum int, err error) {
	if format != "csv" && format != "parquet" {
		return "", "", fmt.Errorf("unsupported format: %s", format)
	}

	// get all battery states of the station
	var states []BatteryState
	if err = s.db.Select(&states, `
		SELECT * FROM battery_state
		WHERE station = ?
		ORDER BY container, pack ASC
	`, station); err != nil {
		return "", "", fmt.Errorf("failed to get battery states: %w", err)
	}

	// generate file
	switch format {
	case "csv":
		file, checksum, err = s.generateCsvFile(states)
	case "parquet":
		file, checksum, err = s.generateParquetFile(states)
	}
	if err != nil {
		return "", "", fmt.Errorf("failed to generate file: %w", err)
	}
}

// generateCsvFile generates csv file of the battery state.
func (s *SqliteStore) generateCsvFile(states []BatteryState) (string, error) {
	if len(states) == 0 {
		return "", nil
	}

	// open a local file for writing, the file name is the timestamp of the snapshot
	t := time.Now()
	file, err := os.Create(fmt.Sprintf("%d/%s/%d.csv", states[0].station, fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day()), t.Unix())
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// create a writebuffer
	w := bufio.NewWriter(file)
	defer w.Flush()

	// crc32 checksum
	crc32 := crc32.NewIEEE()

	// write csv header, delimiter is comma
	// station, container, pack, cell, voltage, current, soc, temperature, state, timestamp
	header := "station,container,pack,cell,voltage,current,soc,temperature,state,timestamp\n"
	w.WriteString(header)
	crc32.Write([]byte(header))

	// write csv body
	for _, state := range states {
		row := fmt.Sprintf("%d,%d,%d,%d,%f,%f,%f,%f,%d\n",
			state.Station, state.Container, state.Pack, state.Cell,
			state.Voltage, state.Current, state.SOC, state.Temperature, state.State, state.Timestamp)
		crc32.Write([]byte(row))
		w.WriteString(row)
	}

	return file.Name(), nil
}

// generateParquetFile generates parquet file of the battery state.
func (s *SqliteStore) generateParquetFile(states []BatteryState) (string, int, error) {
	// TODO
	return "", 0, fmt.Errorf("not implemented")
}
