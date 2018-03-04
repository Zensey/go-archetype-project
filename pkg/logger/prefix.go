package logger

import "time"

const dateLayout = "Jan _2 15:04:05.000"

func getPrefix(tag string, lev LogLevel) string {
	return time.Now().Format(dateLayout) + " [" + tag + "] " + logLevels[lev] + " >"
}

func getPrefixB(tag string, lev LogLevel) string {
	return "[" + tag + "] " + logLevels[lev] + " >"
}