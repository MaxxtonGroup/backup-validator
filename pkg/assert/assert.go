package assert

import (
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type Timings struct {
	RestoreTime time.Duration
	ImportTime  time.Duration
}

type Assert interface {
	RunFor(assertConfig *AssertConfig) bool

	Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string
}
