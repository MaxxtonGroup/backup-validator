package backup

import "time"

type BackupProvider interface {
	Restore(dir string) error

	ListSnapshots(dir string) ([]*Snapshot, error)
}

type Snapshot struct {
	Time time.Time
}
