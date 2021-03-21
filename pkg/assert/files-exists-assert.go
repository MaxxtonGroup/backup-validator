package assert

import (
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
	invalidGlobPatterns := make([]string, 0)

	for _, file := range *assertConfig.FilesExists {
		pattern := filepath.Join(dir, "workdir", file)
		matchingFiles, err := filepath.Glob(pattern)
		if err != nil {
			invalidGlobPatterns = append(invalidGlobPatterns, pattern)
			missingFiles = append(missingFiles, file)
		}

		if len(matchingFiles) == 0 {
			missingFiles = append(missingFiles, file)
		}
	}

	if len(invalidGlobPatterns) > 0 {
		msg := "Invalid Glob patterns: " + strings.Join(invalidGlobPatterns, ", ")
		return &msg
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
