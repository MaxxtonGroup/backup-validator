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
	Name          string
	Duration      time.Duration
	Error         error
	FailedAsserts []string
}

var asserts = []assert.Assert{
	assert.NewFilesExistsAssert(),
	assert.NewFileModifiedAssert(),
	assert.NewBackupRetentionAssert(),
	assert.NewDatabasesExistsAssert(),
	assert.NewDatabasesSizeAssert(),
	assert.NewTablesExistsAssert(),
}

// Validate backups based on tests specified in the configFiles
func Validate(configFiles []string) ([]*TestResult, error) {
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
				log.Printf("Validate backup: %s (running)\n", test.Name)
				startTime := time.Now()
				result, err := validateBackup(&test)
				result.Duration = time.Since(startTime)

				// Collect result
				if err != nil {
					result.Error = err
					log.Printf("Validate backup: %s (failed)\n", test.Name)
				} else {
					log.Printf("Validate backup: %s (done)\n", test.Name)
				}
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func validateBackup(test *TestConfig) (*TestResult, error) {
	result := &TestResult{
		Name: test.Name,
	}

	// create workdir
	dir, err := ioutil.TempDir(".", ".backup-validator")
	if err != nil {
		return result, err
	}
	defer os.RemoveAll(dir)

	// Find runtime provider
	runtimeProvider, err := getRuntimeProvider(test)
	if err != nil {
		return result, err
	}

	// Find backup provider
	backupProvider, err := getBackupProvider(test)
	if err != nil {
		return result, err
	}

	// Find format provider
	formatProvider, err := getFormatProvider(test.Format, runtimeProvider)
	if err != nil {
		return result, err
	}

	// Setup format provider
	err = formatProvider.Setup(dir)
	if err != nil {
		return result, err
	}

	// Destory defer
	defer formatProvider.Destroy(dir)

	// Restore backup
	err = backupProvider.Restore(dir)
	if err != nil {
		return result, err
	}

	// Import backup data in format provider
	err = formatProvider.ImportData(dir, *test.ImportOptions)
	if err != nil {
		return result, err
	}

	// Validate
	if test.Asserts != nil {
		failedAsserts := []string{}
		for _, assertConfig := range *test.Asserts {
			for _, assert := range asserts {
				if assert.RunFor(&assertConfig) {
					msg := assert.Run(dir, &assertConfig, backupProvider, formatProvider)
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

func getFormatProvider(formatType string, runtimeProvider runtime.RuntimeProvider) (format.FormatProvider, error) {
	switch formatType {
	case "file":
		formatProvider := format.NewFileFormatProvider()
		return formatProvider, nil
	case "mongo":
		formatProvider := format.NewMongoFormatProvider(runtimeProvider)
		return formatProvider, nil
	}
	return nil, fmt.Errorf("Unsupported format '%s'", formatType)
}

func getBackupProvider(test *TestConfig) (backup.BackupProvider, error) {
	if test.Restic != nil {
		backupProvider := backup.NewResticBackupProvider(*test.Restic)
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
