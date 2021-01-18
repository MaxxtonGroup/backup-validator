package validator

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

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
				log.Printf("Validation test: %s (running)\n", test.Name)
				result, err := validateBackup(&test)

				// Collect result
				if err != nil {
					result = &TestResult{
						Error: err,
					}
				}
				log.Printf("Validation test: %s (done)\n", test.Name)
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func validateBackup(test *TestConfig) (*TestResult, error) {
	startTime := time.Now()
	result := &TestResult{
		Name: test.Name,
	}

	// Find backup provider
	backupProvider, err := getBackupProvider(test)
	if err != nil {
		return nil, err
	}

	// Find format provider
	formatProvider, err := getFormatProvider(test.Format)
	if err != nil {
		return nil, err
	}

	// Setup format provider
	err = formatProvider.Setup()
	if err != nil {
		return nil, err
	}

	// Restore backup
	err = backupProvider.Restore()
	if err != nil {
		return nil, err
	}

	// Import backup data in format provider
	err = formatProvider.ImportData()
	if err != nil {
		return nil, err
	}

	// Validate
	if test.Asserts != nil {
		failedAsserts := []string{}
		for _, assert := range *test.Asserts {
			msg := formatProvider.Assert(assert)
			if msg != nil {
				failedAsserts = append(failedAsserts, *msg)
			}
		}
		result.FailedAsserts = failedAsserts
	}

	// Destory
	err = formatProvider.Destroy()
	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(startTime)
	return result, nil
}

func getFormatProvider(formatType string) (format.FormatProvider, error) {
	switch formatType {
	case "file":
		formatProvider := format.NewFileFormatProvider()
		return formatProvider, nil
	}
	return nil, fmt.Errorf("Unsupported format '%s'", formatType)
}

func getBackupProvider(test *TestConfig) (backup.BackupProvider, error) {
	if test.Restic != nil {
		backupProvider := backup.NewResticBackupProvider()
		return backupProvider, nil
	}
	return nil, fmt.Errorf("No backup config found")
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
