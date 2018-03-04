// +build !windows

package logger

import (
	"log/syslog"
)

type SyslogBackend struct {
	slog *syslog.Writer
}

func newSyslogBackend(lev LogLevel, tag string) SyslogBackend {
	b := SyslogBackend{}
	var priority syslog.Priority

	switch lev {
	case LogLevelError:
		priority = syslog.LOG_ERR
	case LogLevelWarning:
		priority = syslog.LOG_WARNING
	case LogLevelInfo:
		priority = syslog.LOG_INFO
	case LogLevelDebug:
	case LogLevelTrace:
		priority = syslog.LOG_DEBUG
	}

	b.slog, _ = syslog.New(priority, tag)
	return b
}
func (b SyslogBackend) Write(lev LogLevel, tag, l string) {
	switch lev {
	case LogLevelError:
		b.slog.Err(l)
	case LogLevelWarning:
		b.slog.Warning(l)
	case LogLevelInfo:
		b.slog.Info(l)
	case LogLevelDebug:
	case LogLevelTrace:
		b.slog.Debug(l)
	}
}
