package logger

//go:generate stringer -type=LogLevelRune
type LogLevelRune int

const (
	L_EMERG LogLevelRune = iota + 0
	L_ALERT
	L_CRIT
	L_ERR
	L_WARNING
	L_NOTICE
	L_INFO
	L_DEBUG
)
