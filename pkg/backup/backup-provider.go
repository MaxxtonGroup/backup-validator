package backup

import "time"

type BackupProvider interface {
	Restore(testName string, dir string, snapshot *Snapshot, importOptions []string) error

	ListSnapshots(testName string, dir string) ([]*Snapshot, error)
}

type Snapshot struct {
	Name      string
	Time      time.Time
	Databases []string
}
