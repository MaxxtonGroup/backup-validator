package runtime

type RuntimeProvider interface {
	Setup(testName string, dir string) error
	Destroy(testName string, dir string) error
	Exec(testName string, command string, args ...string) (*string, error)
	ExecRoot(testName string, command string, args ...string) (*string, error)
}
