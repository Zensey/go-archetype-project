package logger

import (
	"fmt"
	"log/syslog"
	"time"
)

//go:generate stringer -type=Severity
type Severity Priority

const (
	L_EMERG   Severity = Severity(LOG_EMERG)
	L_ALERT   Severity = Severity(LOG_ALERT)
	L_CRIT    Severity = Severity(LOG_CRIT)
	L_ERR     Severity = Severity(LOG_ERR)
	L_WARNING Severity = Severity(LOG_WARNING)
	L_NOTICE  Severity = Severity(LOG_NOTICE)
	L_INFO    Severity = Severity(LOG_INFO)
	L_DEBUG   Severity = Severity(LOG_DEBUG)
)

//var Level = L_DEBUG // L_INFO
//const f = log.Ltime |log.Lmicroseconds | log.Lshortfile
//const f = 0
//var Log = log.New(os.Stdout, "",  f)

type Logger struct {
	slog *syslog.Writer
}

func NewLogger(priority syslog.Priority, tag string) (l Logger, err error) {
	l.slog, err = syslog.New(priority, tag)
	return
}

func (l *Logger) Info(a ...interface{}) {
	l.LogHelperln(L_INFO, a...)
}
func (l *Logger) Error(a ...interface{}) {
	l.LogHelperln(L_ERR, a...)
}
func (l *Logger) Debug(a ...interface{}) {
	l.LogHelperln(L_DEBUG, a...)
}

func (l *Logger) Info_f(f string, a ...interface{}) {
	s := fmt.Sprintf(f, a...)
	l.LogHelperln(L_INFO, s)
}
func (l *Logger) Error_f(f string, a ...interface{}) {
	s := fmt.Sprintf(f, a...)
	l.LogHelperln(L_ERR, s)
}
func (l *Logger) Debug_f(f string, a ...interface{}) {
	s := fmt.Sprintf(f, a...)
	l.LogHelperln(L_DEBUG, s)
}

func (l *Logger) LogHelperln(sev Severity, a ...interface{}) {
	prefix := []interface{}{time.Now().Format("Jan _2 15:04:05.000") + " | " + sev.String() + " >"}
	a = append(prefix, a...)
	s_ := fmt.Sprintln(a...)
	fmt.Print(s_)

	switch sev {
	case L_INFO:
		l.slog.Info(s_)
	case L_ERR:
		l.slog.Err(s_)
	case L_DEBUG:
		l.slog.Debug(s_)
	}
}
