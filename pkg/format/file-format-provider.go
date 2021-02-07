package format

import "fmt"

type FileFormatProvider struct {
}

func (p FileFormatProvider) Setup(dir string) error {
	// No setup required
	return nil
}

func (p FileFormatProvider) Destroy(dir string) error {
	// No setup required
	return nil
}

func (p FileFormatProvider) ImportData(dir string, options []string) error {
	// No import required
	return nil
}

func (p FileFormatProvider) GetDatabaseSize(database string) (*uint64, error) {
	return nil, fmt.Errorf(`GetDatabaseSize not available for file format`)
}
func (p FileFormatProvider) ListDatabases() ([]string, error) {
	return nil, fmt.Errorf(`ListDatabases not available for file format`)
}
func (p FileFormatProvider) ListTables(database string) ([]string, error) {
	return nil, fmt.Errorf(`ListTables not available for file format`)
}
func (p FileFormatProvider) QueryRecord(database string, query string) (map[string]interface{}, error) {
	return nil, fmt.Errorf(`QueryRecord not available for file format`)
}

func NewFileFormatProvider() FileFormatProvider {
	fileFormatProvider := FileFormatProvider{}
	return fileFormatProvider
}
