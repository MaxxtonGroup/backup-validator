package assert

import (
	"strings"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type TablesExistsAssert struct {
}

func (a TablesExistsAssert) RunFor(assert *AssertConfig) bool {
	return assert.TablesExists != nil
}

func (a TablesExistsAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string {
	var err error
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
		tables, err := formatProvider.ListTables(testName, db)
		if err != nil {
			msg = err.Error()
			continue
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
			msg = "Missing tables: " + strings.Join(missingTables, ", ")
		} else {
			return nil
		}
	}

	if msg != "" {
		return &msg
	}
	return nil
}

func NewTablesExistsAssert() TablesExistsAssert {
	tablesExistsAssert := TablesExistsAssert{}
	return tablesExistsAssert
}
