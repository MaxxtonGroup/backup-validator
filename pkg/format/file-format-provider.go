package format

type FileFormatProvider struct {
}

func (p FileFormatProvider) Setup(dir string) error {
	// No setup required
	return nil
}

func (p FileFormatProvider) Destroy(dir string) error {
	// No setup required
	return nil
}

func (p FileFormatProvider) ImportData(dir string) error {
	// No import required
	return nil
}

func NewFileFormatProvider() FileFormatProvider {
	fileFormatProvider := FileFormatProvider{}
	return fileFormatProvider
}
