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

func (a FileModifiedAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings, snapshot *backup.Snapshot) *string {
	pattern := filepath.Join(dir, "workdir", assertConfig.FileModified.File)

	// Find matching files
	matchingFiles, err := filepath.Glob(pattern)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid glob pattern: %s", pattern)
		return &errMsg
	}

	if len(matchingFiles) == 0 {
		errMsg := fmt.Sprintf("No matching files for %s", pattern)
		return &errMsg
	}

	// get max duration
	duration, err := time.ParseDuration(assertConfig.FileModified.NewerThan)
	if err != nil {
		errMsg := err.Error()
		return &errMsg
	}

	// Find file within the max duration
	var newestFile string
	var newestMTime *time.Duration = nil
	for _, file := range matchingFiles {
		stats, err := os.Stat(file)
		if err != nil {
			errMsg := err.Error()
			return &errMsg
		}

		mTime := time.Since(stats.ModTime())
		if mTime <= duration {
			return nil
		}
		if newestMTime == nil || mTime < *newestMTime {
			newestFile = file
			newestMTime = &mTime
		}
	}

	msg := fmt.Sprintf("%s is modified %s ago", newestFile, newestMTime)
	return &msg
}

func NewFileModifiedAssert() FileModifiedAssert {
	fileModifiedAssert := FileModifiedAssert{}
	return fileModifiedAssert
}
