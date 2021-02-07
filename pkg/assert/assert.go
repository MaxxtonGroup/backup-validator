package assert

import (
	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type Assert interface {
	RunFor(assertConfig *AssertConfig) bool

	Run(dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider) *string
}
