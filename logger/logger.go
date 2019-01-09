package logger

import (
	"log"
	"os"
)

var infoLogger *log.Logger
var errorLogger *log.Logger

func init() {
	infoLogger = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
}

// Info message into stdout.
func Info(msg string, v ...interface{}) *log.Logger {
	infoLogger.Printf(msg, v...)
	return infoLogger
}

// Error message into stderr.
func Error(msg string, v ...interface{}) *log.Logger {
	errorLogger.Printf(msg, v...)
	return errorLogger
}
