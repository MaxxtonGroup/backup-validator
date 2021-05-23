package assert

import (
	"fmt"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type MaxRestoreTimeAssert struct {
}

func (a MaxRestoreTimeAssert) RunFor(assert *AssertConfig) bool {
	return assert.MaxRestoreTime != nil
}

func (a MaxRestoreTimeAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string {
	maxRestoreTime, err := time.ParseDuration(*assertConfig.MaxRestoreTime)
	if err != nil {
		errMsg := err.Error()
		return &errMsg
	}

	if timings.RestoreTime > maxRestoreTime {
		errMsg := fmt.Sprintf("Restore took %s, which is more than %s", timings.RestoreTime.Round(time.Second), maxRestoreTime.Round(time.Second))
		return &errMsg
	}
	return nil
}

func NewMaxRestoreTimeAssert() MaxRestoreTimeAssert {
	maxRestoreTimeAssert := MaxRestoreTimeAssert{}
	return maxRestoreTimeAssert
}
