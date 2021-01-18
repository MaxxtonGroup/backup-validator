package backup

type ResticConfig struct {
	URL          string `yaml:"url"`
	PasswordFile string `yaml:"passwordFile"`
}
