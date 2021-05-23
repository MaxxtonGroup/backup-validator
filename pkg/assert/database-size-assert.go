package assert

import (
	"fmt"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
	"github.com/dustin/go-humanize"
)

type DatabasesSizeAssert struct {
}

func (a DatabasesSizeAssert) RunFor(assert *AssertConfig) bool {
	return assert.DatabaseSize != nil
}

func (a DatabasesSizeAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string {
	databases, err := formatProvider.ListDatabases(testName)
	if err != nil {
		msg := err.Error()
		return &msg
	}

	databaseName := assertConfig.DatabaseSize.Database
	matchingDatabases, err := getMatchingDatabases(testName, databaseName, databases, formatProvider, *snapshot)
	if err != nil {
		msg := err.Error()
		return &msg
	}

	var msg string
	for _, db := range matchingDatabases {
		size, err := formatProvider.GetDatabaseSize(testName, db)
		if err != nil {
			msg = err.Error()
			continue
		}

		maxSize, err := humanize.ParseBytes(assertConfig.DatabaseSize.Size)
		if err != nil {
			msg = err.Error()
			continue
		}
		if *size < maxSize {
			currentSize := humanize.Bytes(*size)
			msg = fmt.Sprintf("Database size %s is %s, but should be at least %s", db, currentSize, assertConfig.DatabaseSize.Size)
		} else {
			return nil
		}
	}

	if msg != "" {
		return &msg
	}
	return nil
}

func NewDatabasesSizeAssert() DatabasesSizeAssert {
	databasesSizeAssert := DatabasesSizeAssert{}
	return databasesSizeAssert
}
