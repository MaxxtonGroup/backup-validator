package format

type FormatProvider interface {
	Setup() error
	Destroy() error
	ImportData() error
	Assert(assert Assert) *string
}
