package assert

import "github.com/MaxxtonGroup/backup-validator/pkg/backup"

type Assert interface {
	RunFor(assertConfig *AssertConfig) bool

	Run(dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider) *string
}
