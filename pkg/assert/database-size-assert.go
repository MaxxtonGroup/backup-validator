package assert

import (
	"fmt"
	"time"

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
	databaseName := assertConfig.DatabaseSize.Database
	_, isElasticSearch := formatProvider.(format.ElasticsearchFormatProvider)
	if isElasticSearch {
		// Parse database name for elasticsearch
		databaseName = snapshot.Time.Add(-(24 * time.Hour)).Format(databaseName)
	}

	size, err := formatProvider.GetDatabaseSize(testName, databaseName)
	if err != nil {
		msg := err.Error()
		return &msg
	}

	maxSize, err := humanize.ParseBytes(assertConfig.DatabaseSize.Size)
	if err != nil {
		msg := err.Error()
		return &msg
	}
	if *size < maxSize {
		currentSize := humanize.Bytes(*size)
		msg := fmt.Sprintf("Database size %s is %s, but should be at least %s", databaseName, currentSize, assertConfig.DatabaseSize.Size)
		return &msg
	}

	return nil
}

func NewDatabasesSizeAssert() DatabasesSizeAssert {
	databasesSizeAssert := DatabasesSizeAssert{}
	return databasesSizeAssert
}
