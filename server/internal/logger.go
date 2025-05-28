package internal

import (
	"fmt"

	"github.com/fatih/color"
)


const (
	INFO = iota
	WARNING
	ERROR
	DEBUG
)

var logTypeLabels = []string{
	"INFO",
	"WARNING",
	"ERROR",
	"DEBUG",
}

var prefix = "[Chronos] - "
var prefixMigrate = "[Chronos-Migration] - "

func Log(message string, logType int) {
	var logText string

	if logType >= 0 && logType < len(logTypeLabels) {
		logText = prefix + logTypeLabels[logType] + " - " + message
	} else {
		logText = prefix + "SYSTEM - " + message
	}

	switch logType {
		case INFO:
			color.Cyan(logText)
		case WARNING:
			color.Yellow(logText)
		case ERROR:
			color.Red(logText)
		case DEBUG:
			color.White(logText)
		default:
			fmt.Println(logText)
	}
}

func LogMigrate(message string, logType int) {
	var logText string

	if logType >= 0 && logType < len(logTypeLabels) {
		logText = prefixMigrate + logTypeLabels[logType] + " - " + message
	} else {
		logText = prefixMigrate + "SYSTEMMESSAGE - " + message
	}

	switch logType {
		case INFO:
			color.Cyan(logText)
		case WARNING:
			color.Yellow(logText)
		case ERROR:
			color.Red(logText)
		case DEBUG:
			color.White(logText)
		default:
			fmt.Println(logText)
	}
}