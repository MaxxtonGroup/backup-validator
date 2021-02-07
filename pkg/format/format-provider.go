package format

type FormatProvider interface {
	Setup(dir string) error
	Destroy(dir string) error
	ImportData(dir string, options []string) error

	ListDatabases() ([]string, error)
	GetDatabaseSize(database string) (*uint64, error)
	ListTables(database string) ([]string, error)
	QueryRecord(database string, query string) (map[string]interface{}, error)
}
