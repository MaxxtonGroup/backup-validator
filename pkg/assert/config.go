package assert

type AssertConfig struct {
	FilesExists     *[]string                   `yaml:"filesExists"`
	FileModified    *FileModifiedAssertConfig   `yaml:"fileModified"`
	BackupRetention *BackupRetentionAssetConfig `yaml:"backupRetention"`
}

type FileModifiedAssertConfig struct {
	File      string `yaml:"file"`
	NewerThan string `yaml:"newerThan"`
}

type BackupRetentionAssetConfig struct {
	Snapshots *int    `yaml:"snapshots"`
	OlderThan *string `yaml:"olderThan"`
}
