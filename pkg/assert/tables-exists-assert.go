package assert

import (
	"strings"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type TablesExistsAssert struct {
}

func (a TablesExistsAssert) RunFor(assert *AssertConfig) bool {
	return assert.TablesExists != nil
}

func (a TablesExistsAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string {
	databaseName := assertConfig.TablesExists.Database
	_, isElasticSearch := formatProvider.(format.ElasticsearchFormatProvider)
	if isElasticSearch {
		// Parse database name for elasticsearch
		databaseName = snapshot.Time.Add(-(24 * time.Hour)).Format(databaseName)
	}
	tables, err := formatProvider.ListTables(testName, databaseName)
	if err != nil {
		msg := err.Error()
		return &msg
	}

	missingTables := make([]string, 0)
	for _, tableName := range *assertConfig.TablesExists.Tables {
		exists := false
		for _, table := range tables {
			if table == tableName {
				exists = true
			}
		}
		if !exists {
			missingTables = append(missingTables, tableName)
		}
	}

	if len(missingTables) > 0 {
		msg := "Missing tables: " + strings.Join(missingTables, ", ")
		return &msg
	}
	return nil
}

func NewTablesExistsAssert() TablesExistsAssert {
	tablesExistsAssert := TablesExistsAssert{}
	return tablesExistsAssert
}
