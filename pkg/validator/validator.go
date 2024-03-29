package validator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/assert"
	"github.com/MaxxtonGroup/backup-validator/pkg/runtime"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"

	"github.com/MaxxtonGroup/backup-validator/pkg/format"
	"github.com/ghodss/yaml"
)

type TestResult struct {
	Name            string        `json:"name"`
	TotalDuration   time.Duration `json:"totalDuration"`
	RestoreDuration time.Duration `json:"restoreDuration"`
	ImportDuration  time.Duration `json:"importDuration"`
	Error           *string       `json:"error"`
	FailedAsserts   []string      `json:"failedAsserts"`
}

var asserts = []assert.Assert{
	assert.NewFilesExistsAssert(),
	assert.NewFileModifiedAssert(),
	assert.NewBackupRetentionAssert(),
	assert.NewMaxRestoreTimeAssert(),
	assert.NewMaxImportTimeAssert(),
	assert.NewDatabasesExistsAssert(),
	assert.NewDatabasesSizeAssert(),
	assert.NewTablesExistsAssert(),
}

// Validate backups based on tests specified in the configFiles
func Validate(configFiles []string, cleanup bool) ([]*TestResult, error) {
	// Load config files
	configs, err := loadConfig(configFiles)
	if err != nil {
		return nil, err
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("No config files provided, use --config-file=<file> to provide one")
	}

	// Run tests in serial
	log.Println("Starting Test Suite")
	results := make([]*TestResult, 0)
	for _, config := range configs {
		if *config.Tests != nil {
			for _, test := range *config.Tests {

				// Run test
				log.Printf("[%s] Validate backup (running)\n", test.Name)
				startTime := time.Now()
				result, err := validateBackup(&test, cleanup)
				result.TotalDuration = time.Since(startTime)

				// Collect result
				if err != nil {
					errMsg := err.Error()
					result.Error = &errMsg
					log.Printf("[%s] Validate backup (failed)\n", test.Name)
				} else {
					log.Printf("[%s] Validate backup (done)\n", test.Name)
				}
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func validateBackup(test *TestConfig, cleanup bool) (*TestResult, error) {
	result := &TestResult{
		Name: test.Name,
	}

	// create workdir
	dir, err := ioutil.TempDir(".", ".backup-validator")
	if err != nil {
		return result, err
	}
	if cleanup {
		defer os.RemoveAll(dir)
	}

	// Find runtime provider
	runtimeProvider, err := getRuntimeProvider(test)
	if err != nil {
		return result, err
	}

	// Find backup provider
	backupProvider, err := getBackupProvider(test, runtimeProvider)
	if err != nil {
		return result, err
	}

	// Find format provider
	formatProvider, err := getFormatProvider(test.Format, runtimeProvider, test)
	if err != nil {
		return result, err
	}

	// Destory format provider
	if cleanup {
		defer formatProvider.Destroy(test.Name, dir)
	}

	// Setup format provider
	for i := 0; i < 5; i++ {
		err = formatProvider.Setup(test.Name, dir)
		if err == nil {
			break
		} else {
			log.Printf("[%s] Setup failed %s, retrying...", test.Name, err)
			if i >= 4 {
				return result, err
			}
		}
	}

	// Find last snapshot
	snapshots, err := backupProvider.ListSnapshots(test.Name, dir)
	if err != nil {
		return result, err
	}
	if len(snapshots) == 0 {
		return result, fmt.Errorf("no snapshots found")
	}
	snapshot := snapshots[len(snapshots)-1]

	if test.ImportOptions == nil {
		importOptions := []string{}
		test.ImportOptions = &importOptions
	}

	// Restore backup
	restoreStartTime := time.Now()
	err = backupProvider.Restore(test.Name, dir, snapshot, *test.ImportOptions)
	result.RestoreDuration = time.Since(restoreStartTime)
	if err != nil {
		return result, err
	}

	// Import backup data in format provider
	log.Printf("[%s] Importing data...\n", test.Name)
	importStartTime := time.Now()
	err = formatProvider.ImportData(test.Name, dir, *test.ImportOptions)
	result.ImportDuration = time.Since(importStartTime)
	if err != nil {
		return result, err
	}

	// Validate
	if test.Asserts != nil {
		timings := assert.Timings{
			RestoreTime: result.RestoreDuration,
			ImportTime:  result.ImportDuration,
		}
		failedAsserts := []string{}
		for _, assertConfig := range *test.Asserts {
			for _, assert := range asserts {
				if assert.RunFor(&assertConfig) {
					msg := assert.Run(test.Name, dir, &assertConfig, backupProvider, formatProvider, timings, snapshot)
					if msg != nil {
						failedAsserts = append(failedAsserts, *msg)
					}
				}
			}
		}
		result.FailedAsserts = failedAsserts
	}

	return result, nil
}

func getFormatProvider(formatType string, runtimeProvider runtime.RuntimeProvider, test *TestConfig) (format.FormatProvider, error) {
	switch formatType {
	case "file":
		formatProvider := format.NewFileFormatProvider()
		return formatProvider, nil
	case "mongo":
		formatProvider := format.NewMongoFormatProvider(runtimeProvider)
		return formatProvider, nil
	case "postgresql":
		formatProvider := format.NewPostgresqlFormatProvider(runtimeProvider)
		return formatProvider, nil
	case "elasticsearch":
		formatProvider := format.NewElasticsearchFormatProvider(runtimeProvider, *test.ElasticsearchSnapshotRepository)
		return formatProvider, nil
	}
	return nil, fmt.Errorf("Unsupported format '%s'", formatType)
}

func getBackupProvider(test *TestConfig, runtimeProvider runtime.RuntimeProvider) (backup.BackupProvider, error) {
	if test.Restic != nil {
		backupProvider := backup.NewResticBackupProvider(*test.Restic)
		return backupProvider, nil
	}
	if test.ElasticsearchSnapshotRepository != nil {
		backupProvider := backup.NewElasticsearchBackupProvider(runtimeProvider)
		return backupProvider, nil
	}
	return nil, fmt.Errorf("No backup config found")
}

func getRuntimeProvider(test *TestConfig) (runtime.RuntimeProvider, error) {
	if test.Docker != nil {
		runtimeProvider := runtime.NewDockerRuntimeProvider(*test.Docker)
		return runtimeProvider, nil
	}
	return nil, nil
}

func loadConfig(configFiles []string) ([]*ValidatorConfig, error) {
	configs := make([]*ValidatorConfig, 0)
	for _, configFile := range configFiles {
		yamlRaw, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
		config := &ValidatorConfig{}
		err = yaml.Unmarshal(yamlRaw, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, nil
}
