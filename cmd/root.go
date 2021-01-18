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

	"github.com/MaxxtonGroup/backup-validator/pkg/validator"

	"github.com/spf13/cobra"
)

var configFiles []string = []string{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backup-validator",
	Short: "CLI to validate backups by restoring them",
	Long:  `backup-validator is a CLI for validating Restic/Elasticsearch backups by restoring them`,
	Run: func(cmd *cobra.Command, args []string) {
		// Execute command
		testResults, err := validator.Validate(configFiles)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		log.Println("")
		log.Println("Test result:")
		for _, testResult := range testResults {
			log.Printf("- %s (%s):", testResult.Name, testResult.Duration)
			if testResult.Error != nil {
				log.Printf("\n    error: %s\n", testResult.Error)
			} else if testResult.FailedAsserts != nil && len(testResult.FailedAsserts) > 0 {
				for _, failedAssert := range testResult.FailedAsserts {
					log.Printf("\n    assert failed: %s\n", failedAssert)
				}
			} else {
				log.Println("    successful")
			}
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
	rootCmd.Flags().StringSliceVar(&configFiles, "config-file", []string{}, "Config files")
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
