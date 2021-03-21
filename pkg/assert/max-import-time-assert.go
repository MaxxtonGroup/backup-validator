package assert

import (
	"fmt"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type MaxImportTimeAssert struct {
}

func (a MaxImportTimeAssert) RunFor(assert *AssertConfig) bool {
	return assert.MaxImportTime != nil
}

func (a MaxImportTimeAssert) Run(dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings) *string {
	maxImportTime, err := time.ParseDuration(*assertConfig.MaxImportTime)
	if err != nil {
		errMsg := err.Error()
		return &errMsg
	}

	if timings.ImportTime > maxImportTime {
		errMsg := fmt.Sprintf("Importing database took %s, which is more than %s", timings.ImportTime.Round(time.Second), maxImportTime.Round(time.Second))
		return &errMsg
	}
	return nil
}

func NewMaxImportTimeAssert() MaxImportTimeAssert {
	maxImportTimeAssert := MaxImportTimeAssert{}
	return maxImportTimeAssert
}
