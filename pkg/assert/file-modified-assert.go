package assert

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type FileModifiedAssert struct {
}

func (a FileModifiedAssert) RunFor(assert *AssertConfig) bool {
	return assert.FileModified != nil
}

func (a FileModifiedAssert) Run(dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider) *string {
	filePath := filepath.Join(dir, "workdir", assertConfig.FileModified.File)
	duration, err := time.ParseDuration(assertConfig.FileModified.NewerThan)
	if err != nil {
		errMsg := err.Error()
		return &errMsg
	}

	stats, err := os.Stat(filePath)
	if err != nil {
		errMsg := err.Error()
		return &errMsg
	}

	mTime := time.Since(stats.ModTime())
	if mTime > duration {
		msg := fmt.Sprintf("%s is modified %s ago", assertConfig.FileModified.File, mTime)
		return &msg
	}

	return nil
}

func NewFileModifiedAssert() FileModifiedAssert {
	fileModifiedAssert := FileModifiedAssert{}
	return fileModifiedAssert
}
