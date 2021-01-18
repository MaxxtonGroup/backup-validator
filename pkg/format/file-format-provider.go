package format

type FileFormatProvider struct {
}

func (p FileFormatProvider) Setup() error {
	// No setup required
	return nil
}

func (p FileFormatProvider) Destroy() error {
	// No setup required
	return nil
}

func (p FileFormatProvider) ImportData() error {
	// No import required
	return nil
}

func (p FileFormatProvider) Assert(assert Assert) *string {
	// No import required
	return nil
}

func NewFileFormatProvider() FileFormatProvider {
	fileFormatProvider := FileFormatProvider{}
	return fileFormatProvider
}
