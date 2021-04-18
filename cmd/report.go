package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/MaxxtonGroup/backup-validator/pkg/report"
	"github.com/MaxxtonGroup/backup-validator/pkg/validator"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var reportCmd = &cobra.Command{
	Use:   "report <report_files>",
	Short: "Generate report based on one ore more report.json file(s)",
	Long:  `Generate report based on one ore more report.json file(s)`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		testResults := make([]*validator.TestResult, 0)

		// Validate options
		if reportFormat != "" && reportFormat != "json" && reportFormat != "html" {
			fmt.Println("Error: invalid flag: --report-format should be one of: \"json\" or \"html\"")
			os.Exit(1)
		}

		// read files
		for _, reportFile := range args {
			err, results := report.LoadJsonReport(reportFile)
			if err != nil {
				log.Printf("Failed to load json report '%s': %s", reportFile, err)
				os.Exit(1)
			}
			testResults = append(testResults, results...)
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
		} else {
			log.Printf("report-file is empty\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.Flags().StringVarP(&reportFile, "report-file", "o", "report.json", "Output file for the test results.")
	reportCmd.Flags().StringVarP(&reportFormat, "report-format", "", "json", "Format of the test results. One of: \"json\" or \"html\".")
}
