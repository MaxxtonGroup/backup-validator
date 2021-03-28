package assert

import (
	"fmt"
	"time"

	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
)

type BackupRetentionAssert struct {
}

func (a BackupRetentionAssert) RunFor(assert *AssertConfig) bool {
	return assert.BackupRetention != nil
}

func (a BackupRetentionAssert) Run(testName string, dir string, assertConfig *AssertConfig, backupProvider backup.BackupProvider, formatProvider format.FormatProvider, timings Timings) *string {
	snapshots, err := backupProvider.ListSnapshots(testName, dir)
	if err != nil {
		msg := err.Error()
		return &msg
	}

	if assertConfig.BackupRetention.Snapshots != nil {
		if len(snapshots) < *assertConfig.BackupRetention.Snapshots {
			msg := fmt.Sprintf("There are only %d snapshots available", len(snapshots))
			return &msg
		}
	}

	if assertConfig.BackupRetention.OlderThan != nil {
		duration, err := time.ParseDuration(*assertConfig.BackupRetention.OlderThan)
		if err != nil {
			msg := err.Error()
			return &msg
		}

		var oldestSnapshot *time.Time
		for _, snapshot := range snapshots {
			if oldestSnapshot == nil || oldestSnapshot.After(snapshot.Time) {
				oldestSnapshot = &snapshot.Time
			}
		}
		if oldestSnapshot == nil {
			msg := "No snapshots found"
			return &msg
		} else {
			diff := time.Since(*oldestSnapshot)
			if diff < duration {
				msg := fmt.Sprintf("Oldest snapshot is from %s ago", diff)
				return &msg
			}
		}
	}

	return nil
}

func NewBackupRetentionAssert() BackupRetentionAssert {
	backupRetentionAssert := BackupRetentionAssert{}
	return backupRetentionAssert
}
