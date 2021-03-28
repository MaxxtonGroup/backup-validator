package format

type FormatProvider interface {
	Setup(testName string, dir string) error
	Destroy(testName string, dir string) error
	ImportData(testName string, dir string, options []string) error

	ListDatabases(testName string) ([]string, error)
	GetDatabaseSize(testName string, database string) (*uint64, error)
	ListTables(testName string, database string) ([]string, error)
	QueryRecord(testName string, database string, query string) (map[string]interface{}, error)
}
