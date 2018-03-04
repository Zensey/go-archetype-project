// +build windows

package logger

type SyslogBackend struct {
}

func newSyslogBackend(lev LogLevel, tag string) SyslogBackend {
	return SyslogBackend{}
}

func (b SyslogBackend) Write(lev LogLevel, tag, l string) {}
