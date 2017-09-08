package logger

import (
	"fmt"
	"log"
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


var Level = L_DEBUG // L_INFO
const f = log.Ltime |log.Lmicroseconds | log.Lshortfile
//const f = 0
//var Log = log.New(os.Stdout, "",  f)


func Log_info (a ...interface{}) {
	LogHelperln(L_INFO, a...)
}
func Log_debug (a ...interface{}) {
	LogHelperln(L_DEBUG, a...)
}
func Log_err (a ...interface{}) {
	LogHelperln(L_ERR, a...)
}

func Log_info_f(f string, a ...interface{}) {
	LogHelper_f(L_INFO, f, a...)
}
func Log_debug_f(f string, a ...interface{}) {
	LogHelper_f(L_DEBUG, f, a...)
}
func Log_err_f (f string, a ...interface{}) {
	LogHelper_f(L_ERR, f, a...)
}

func LogHelperln(sev Severity, a ...interface{}) {
	if sev <= Level {
		s := []interface{}{ sev.String()+" >" }
		s = append(s, a...)
		fmt.Println(s...)
	}
}

func LogHelper_f(sev Severity, f string, a ...interface{}) {
	if sev <= Level {
		//fmt.Printf(sev.String()+" > " + f, a...)
		log.Printf(sev.String()+" > " + f, a...)
	}
}
