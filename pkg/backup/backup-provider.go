package backup

import "time"

type BackupProvider interface {
	Restore(testName string, dir string) error

	ListSnapshots(testName string, dir string) ([]*Snapshot, error)
}

type Snapshot struct {
	Time time.Time
}
