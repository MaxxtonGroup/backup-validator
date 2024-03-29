package assert

type AssertConfig struct {
	FilesExists     *[]string                    `yaml:"filesExists"`
	FileModified    *FileModifiedAssertConfig    `yaml:"fileModified"`
	BackupRetention *BackupRetentionAssertConfig `yaml:"backupRetention"`
	MaxRestoreTime  *string                      `yaml:"maxRestoreTime"`
	MaxImportTime   *string                      `yaml:"maxImportTime"`

	DatabasesExists *[]string                 `yaml:"databasesExists"`
	DatabaseSize    *DatabaseSizeAssertConfig `yaml:"databaseSize"`
	TablesExists    *TableExistsAssertConfig  `yaml:"tablesExists"`
	QueryRecord     *QueryRecordAssertConfig  `yaml:"queryRecord"`
}

type FileModifiedAssertConfig struct {
	File      string `yaml:"file"`
	NewerThan string `yaml:"newerThan"`
}

type BackupRetentionAssertConfig struct {
	Snapshots *int    `yaml:"snapshots"`
	OlderThan *string `yaml:"olderThan"`
}

type QueryRecordAssertConfig struct {
	Database string `yaml:"database"`
	Query    string `yaml:"query"`
	Matches  map[string]interface{}
}

type DatabaseSizeAssertConfig struct {
	Database string `yaml:"database"`
	Size     string `yaml:"size"`
}

type TableExistsAssertConfig struct {
	Database string    `yaml:"database"`
	Tables   *[]string `yaml:"tables"`
}
