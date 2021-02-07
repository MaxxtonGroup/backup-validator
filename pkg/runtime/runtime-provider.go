package runtime

type RuntimeProvider interface {
	Setup(dir string) error
	Destroy(dir string) error
	Exec(command string, args ...string) (*string, error)
}
