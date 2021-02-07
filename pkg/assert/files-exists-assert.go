package assert

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type FilesExistsAssert struct {
}

func (a FilesExistsAssert) RunFor(assert *AssertConfig) bool {
	return assert.FilesExists != nil
}

func (a FilesExistsAssert) Run(dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider) *string {
	missingFiles := make([]string, 0)

	for _, file := range *assertConfig.FilesExists {
		filePath := filepath.Join(dir, "workdir", file)
		_, err := os.Stat(filePath)
		if err != nil {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(missingFiles) > 0 {
		msg := "Missing files: " + strings.Join(missingFiles, ", ")
		return &msg
	}
	return nil
}

func NewFilesExistsAssert() FilesExistsAssert {
	filesExistsAssert := FilesExistsAssert{}
	return filesExistsAssert
}
