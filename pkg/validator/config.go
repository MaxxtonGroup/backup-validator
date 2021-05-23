package validator

import (
	"github.com/MaxxtonGroup/backup-validator/pkg/assert"
	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
	"github.com/MaxxtonGroup/backup-validator/pkg/format"
	"github.com/MaxxtonGroup/backup-validator/pkg/runtime"
)

type ValidatorConfig struct {
	Tests *[]TestConfig `yaml:"tets"`
}

type TestConfig struct {
	Name   string `yaml:"name"`
	Format string `yaml:"format"`

	Restic                          *backup.ResticConfig                    `yaml:"restic"`
	ElasticsearchSnapshotRepository *format.ElasticsearchSnapshotRepository `yaml:"elasticsearchSnapshotRepository"`
	Asserts                         *[]assert.AssertConfig                  `yaml:"asserts"`
	Docker                          *runtime.DockerConfig                   `yaml:"docker"`
	ImportOptions                   *[]string                               `yaml:"importOptions"`
}

type DockerConfig struct {
	Image string `yaml:"image"`
}
