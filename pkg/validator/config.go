package validator

import (
	"github.com/MaxxtonGroup/backup-validator/pkg/assert"
	"github.com/MaxxtonGroup/backup-validator/pkg/backup"
)

type ValidatorConfig struct {
	Tests *[]TestConfig `yaml:"tets"`
}

type TestConfig struct {
	Name   string `yaml:"name"`
	Format string `yaml:"format"`

	Restic  *backup.ResticConfig `yaml:"restic"`
	Asserts *[]assert.AssertConfig     `yaml:"asserts"`
}

type DockerConfig struct {
	Image string `yaml:"image"`
}
