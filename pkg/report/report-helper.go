package report

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/validator"
)

type TemplateReport struct {
	TestResults []*TemplateTestResult
	Time        string
}

type TemplateTestResult struct {
	Name            string
	TotalDuration   string
	RestoreDuration string
	ImportDuration  string
	Error           *string
	FailedAsserts   []string
}

func StoreJsonReport(reportFile string, testResults []*validator.TestResult) error {
	str, err := json.Marshal(testResults)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(reportFile, str, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func LoadJsonReport(reportFile string) (error, []*validator.TestResult) {
	bytes, err := ioutil.ReadFile(reportFile)
	if err != nil {
		return err, nil
	}

	testResults := make([]*validator.TestResult, 0)
	err = json.Unmarshal(bytes, &testResults)
	if err != nil {
		return err, nil
	}

	return nil, testResults
}

func StoreHtmlReport(reportFile string, testResults []*validator.TestResult) error {
	templateTestResults := make([]*TemplateTestResult, 0)
	for _, result := range testResults {
		templateResult := TemplateTestResult{
			Name:            result.Name,
			TotalDuration:   result.TotalDuration.Round(time.Second).String(),
			RestoreDuration: result.RestoreDuration.Round(time.Second).String(),
			ImportDuration:  result.ImportDuration.Round(time.Second).String(),
			Error:           result.Error,
			FailedAsserts:   result.FailedAsserts,
		}
		templateTestResults = append(templateTestResults, &templateResult)
	}
	report := TemplateReport{
		TestResults: templateTestResults,
		Time:        time.Now().Format(time.RFC3339),
	}

	tmpl, err := template.New("email").Parse(HTML_REPORT)
	if err != nil {
		return err
	}

	f, err := os.Create(reportFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, report)
	if err != nil {
		return err
	}

	return nil
}
