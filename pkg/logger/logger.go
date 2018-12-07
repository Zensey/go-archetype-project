package logger

import (
	"fmt"
)

type LogLevel int

const (
	LogLevelNull LogLevel = iota
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

var logLevels = map[LogLevel]string{
	LogLevelNull:    "",
	LogLevelError:   "error",
	LogLevelWarning: "warning",
	LogLevelInfo:    "info",
	LogLevelDebug:   "debug",
	LogLevelTrace:   "trace",
}

type LoggerBackend interface {
	Write(lev LogLevel, tag, l string)
}

type Logger struct {
	tag   string
	level LogLevel
	back  LoggerBackend
}

//////////////
type BackendID int

const (
	BackendConsole BackendID = iota
	BackendSyslog
	BackendFile
	BackendNull
)

var logBack = map[BackendID]string{
	BackendConsole: "console",
	BackendSyslog:  "syslog",
	BackendFile:    "file",
	BackendNull:    "null",
}

func LookupLogLevel(name string) LogLevel {
	for k, v := range logLevels {
		if v == name {
			return k
		}
	}
	return LogLevelNull
}

func LookupLogBackend(name string) BackendID {
	for k, v := range logBack {
		if v == name {
			return k
		}
	}
	return BackendNull
}

func NewLogger(lev LogLevel, tag string, backend BackendID) (l Logger, err error) {
	l.tag = tag
	l.level = lev

	switch backend {
	case BackendConsole:
		l.back = newConsoleBackend()
	case BackendSyslog:
		l.back, err = newSyslogBackend(lev, tag)
	case BackendFile:
		l.back = newFileBackend(lev, tag)
	case BackendNull:
		l.back = newNullBackend()
	default:
		panic("not supported logging backend")
	}
	return
}

func (l *Logger) Error(a ...interface{}) {
	l.logPrint(LogLevelError, a...)
}
func (l *Logger) Errorf(f string, a ...interface{}) {
	l.logPrintf(LogLevelError, f, a...)
}

func (l *Logger) Warning(a ...interface{}) {
	l.logPrint(LogLevelWarning, a...)
}
func (l *Logger) Warningf(f string, a ...interface{}) {
	l.logPrintf(LogLevelWarning, f, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.logPrint(LogLevelInfo, a...)
}
func (l *Logger) Infof(f string, a ...interface{}) {
	l.logPrintf(LogLevelInfo, f, a...)
}

func (l *Logger) Debug(a ...interface{}) {
	l.logPrint(LogLevelDebug, a...)
}
func (l *Logger) Debugf(f string, a ...interface{}) {
	l.logPrintf(LogLevelDebug, f, a...)
}

func (l *Logger) Trace(a ...interface{}) {
	l.logPrint(LogLevelTrace, a...)
}
func (l *Logger) Tracef(f string, a ...interface{}) {
	l.logPrintf(LogLevelTrace, f, a...)
}

//////

func (l *Logger) logPrintf(sev LogLevel, format string, a ...interface{}) {
	l.toBackEnd(sev, fmt.Sprintf(format, a...))
}

func (l *Logger) logPrint(sev LogLevel, a ...interface{}) {
	s := fmt.Sprintln(a...)
	s = s[0 : len(s)-1]
	l.toBackEnd(sev, s)
}

func (l *Logger) toBackEnd(level LogLevel, s string) {
	if l.level < level {
		return
	}
	l.back.Write(level, l.tag, s)
}
