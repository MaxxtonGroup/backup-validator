package assert

import (
	"strings"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type DatabasesExistsAssert struct {
}

func (a DatabasesExistsAssert) RunFor(assert *AssertConfig) bool {
	return assert.DatabasesExists != nil
}

func (a DatabasesExistsAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string {
	var err error
	databases := snapshot.Databases
	if databases == nil {
		databases, err = formatProvider.ListDatabases(testName)
		if err != nil {
			msg := err.Error()
			return &msg
		}
	}

	missingDatabases := make([]string, 0)
	for _, databaseName := range *assertConfig.DatabasesExists {
		_, isElasticSearch := formatProvider.(format.ElasticsearchFormatProvider)
		if isElasticSearch {
			// Parse database name for elasticsearch
			databaseName = snapshot.Time.Add(-(24 * time.Hour)).Format(databaseName)
		}
		exists := false
		for _, database := range databases {
			if database == databaseName {
				exists = true
			}
		}
		if !exists {
			missingDatabases = append(missingDatabases, databaseName)
		}
	}

	if len(missingDatabases) > 0 {
		msg := "Missing databases: " + strings.Join(missingDatabases, ", ")
		return &msg
	}
	return nil
}

func NewDatabasesExistsAssert() DatabasesExistsAssert {
	databasesExistsAssert := DatabasesExistsAssert{}
	return databasesExistsAssert
}
