package logger

import (
	"log"
	"os"
	"sync"
)

var infoLogger *log.Logger
var errorLogger *log.Logger

var once sync.Once

// Info message into stdout.
func Info(msg string, v ...interface{}) *log.Logger {
	once.Do(func() {
		infoLogger = log.New(os.Stdout, "", 0)
	})

	infoLogger.Printf(msg, v...)

	return infoLogger
}

// Error message into stderr.
func Error(msg string, v ...interface{}) *log.Logger {
	once.Do(func() {
		errorLogger = log.New(os.Stderr, "", 0)
	})

	errorLogger.Printf(msg, v...)

	return errorLogger
}
