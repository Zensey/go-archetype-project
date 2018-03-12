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

func NewLogger(lev LogLevel, tag, back string) (l Logger, err error) {
	l.tag = tag
	l.level = lev

	switch back {
	case "console":
		l.back = newConsoleBackend()
	case "syslog":
		l.back = newSyslogBackend(lev, tag)
	case "file":
		l.back = newFileBackend(lev, tag)
	case "null":
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
	s = s[0:len(s)-1]
	l.toBackEnd(sev, s)
}

func (l *Logger) toBackEnd(level LogLevel, s string) {
	if l.level < level {
		return
	}
	l.back.Write(level, l.tag, s)
}
