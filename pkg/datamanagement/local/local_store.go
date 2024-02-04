package localstore

import (
	"github.com/zhangjinpeng87/openbms/pkg/datamodel"
)

// LocalStore is the interface of local data store.
// It is used to store real-time batteries data locally.
// The real-time batteries data is used to do e.
type LocalStore interface {
	// Open opens the local store.
	Open() error

	// Close closes the local store.
	Close() error

	// Update or insert a battery state.
	Upsert(state *datamodel.BatteryState) error

	// Get the latest battery state.
	GetLatest(station, container, pack, cell int) (*datamodel.BatteryState, error)

	// Generate snapshot file of the battery state for specified station.
	// The snapshot file will be upload to the cloud storage.
	// Format can be "csv" or "parquet".
	// return file name, file checksum and error
	GenerateSnapshotFile(station int, format string) (string, int, error)
}
