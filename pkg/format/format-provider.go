package format

type FormatProvider interface {
	Setup(dir string) error
	Destroy(dir string) error
	ImportData(dir string) error
}
