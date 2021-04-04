/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/report"
	"github.com/MaxxtonGroup/backup-validator/pkg/validator"

	"github.com/spf13/cobra"
)

var configFiles []string = []string{}
var cleanup bool
var reportFile string
var reportFormat string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backup-validator",
	Short: "CLI to validate backups by restoring them",
	Long:  `backup-validator is a CLI for validating Restic/Elasticsearch backups by restoring them`,
	Run: func(cmd *cobra.Command, args []string) {

		// Validate options
		if reportFormat != "" && reportFormat != "json" && reportFormat != "html" {
			fmt.Println("Error: invalid flag: --report-format should be one of: \"json\" or \"html\"")
			os.Exit(1)
		}

		// Execute command
		testResults, err := validator.Validate(configFiles, cleanup)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		failedTests := 0

		log.Println("")
		log.Println("Test result:")
		for _, testResult := range testResults {
			log.Printf("- %s (total: %s, restore: %s, import: %s):", testResult.Name, testResult.TotalDuration.Round(time.Second), testResult.RestoreDuration.Round(time.Second), testResult.ImportDuration.Round(time.Second))
			if testResult.Error != nil {
				failedTests++
				log.Printf("    error: %s\n", *testResult.Error)
			} else if testResult.FailedAsserts != nil && len(testResult.FailedAsserts) > 0 {
				for _, failedAssert := range testResult.FailedAsserts {
					failedTests++
					log.Printf("    assert failed: %s", failedAssert)
				}
				log.Println()
			} else {
				log.Println("    valid")
			}
		}

		// generate reports
		if reportFile != "" {
			if reportFormat == "json" {
				err = report.StoreJsonReport(reportFile, testResults)
				if err != nil {
					log.Printf("Failed to create json report: %s", err)
					os.Exit(1)
				}
			}
			if reportFormat == "html" {
				err = report.StoreHtmlReport(reportFile, testResults)
				if err != nil {
					log.Printf("Failed to create html report: %s", err)
					os.Exit(1)
				}
			}
		}

		if failedTests > 0 {
			os.Exit(1)
		}
	},
}

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.opsgy.yml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringSliceVarP(&configFiles, "test-file", "f", []string{}, "Test definition files.")
	rootCmd.Flags().BoolVarP(&cleanup, "cleanup", "c", true, "Cleanup backup files after test has finished.")
	rootCmd.Flags().StringVarP(&reportFile, "report-file", "o", "report.json", "Output file for the test results.")
	rootCmd.Flags().StringVarP(&reportFormat, "report-format", "", "json", "Format of the test results. One of: \"json\" or \"html\".")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	opsgy.SetConfigFile(cfgFile)
	// } else {
	// 	opsgy.SetConfigFile(opsgy.GetDefaultConfigFile())
	// }
}
