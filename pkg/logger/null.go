package logger

type NullBackend struct {}

func newNullBackend() NullBackend {
	return NullBackend{}
}

func (b NullBackend) Write(lev LogLevel, tag, l string) {
}
