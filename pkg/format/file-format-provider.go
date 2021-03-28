package format

import "fmt"

type FileFormatProvider struct {
}

func (p FileFormatProvider) Setup(testName string, dir string) error {
	// No setup required
	return nil
}

func (p FileFormatProvider) Destroy(testName string, dir string) error {
	// No setup required
	return nil
}

func (p FileFormatProvider) ImportData(testName string, dir string, options []string) error {
	// No import required
	return nil
}

func (p FileFormatProvider) GetDatabaseSize(testName string, database string) (*uint64, error) {
	return nil, fmt.Errorf(`[%s] GetDatabaseSize not available for file format`, testName)
}
func (p FileFormatProvider) ListDatabases(testName string) ([]string, error) {
	return nil, fmt.Errorf(`[%s] ListDatabases not available for file format`, testName)
}
func (p FileFormatProvider) ListTables(testName string, database string) ([]string, error) {
	return nil, fmt.Errorf(`[%s] ListTables not available for file format`, testName)
}
func (p FileFormatProvider) QueryRecord(testName string, database string, query string) (map[string]interface{}, error) {
	return nil, fmt.Errorf(`[%s] QueryRecord not available for file format`, testName)
}

func NewFileFormatProvider() FileFormatProvider {
	fileFormatProvider := FileFormatProvider{}
	return fileFormatProvider
}
