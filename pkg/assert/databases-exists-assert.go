package assert

import (
	"fmt"
	"regexp"
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
		matchingDatabases, err := getMatchingDatabases(testName, databaseName, databases, formatProvider, *snapshot)
		if err != nil {
			msg := err.Error()
			return &msg
		}

		if len(matchingDatabases) == 0 {
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

func getMatchingDatabases(testName string, databaseName string, databases []string, formatProvider format.FormatProvider, snapshot backup.Snapshot) ([]string, error) {
	_, isElasticSearch := formatProvider.(format.ElasticsearchFormatProvider)
	if isElasticSearch {
		// Parse database name for elasticsearch
		databaseName = snapshot.Time.Add(-(24 * time.Hour)).Format(databaseName)
	}
	parts := strings.Split(databaseName, "*")
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}
	databaseRegexStr := strings.Join(parts, ".*")
	databaseRegex, err := regexp.Compile(databaseRegexStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid regex: " + databaseRegexStr)
	}

	matchingDatabases := []string{}
	for _, database := range databases {
		if databaseRegex.MatchString(database) {
			matchingDatabases = append(matchingDatabases, database)
		}
	}
	return matchingDatabases, nil
}
