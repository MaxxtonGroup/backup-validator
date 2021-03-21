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

func (a TablesExistsAssert) Run(dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings) *string {
	tables, err := formatProvider.ListTables(assertConfig.TablesExists.Database)
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
