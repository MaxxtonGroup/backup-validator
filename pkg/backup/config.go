package backup

type ResticConfig struct {
	Repository   string            `yaml:"repository"`
	PasswordFile string            `yaml:"passwordFile"`
	Password     *string           `yaml:"password"`
	Env          map[string]string `yaml:"env"`
}
