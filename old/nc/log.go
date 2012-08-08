package nc

// Logging and error reporting utility functions

import (
	"fmt"
	"log"
)

var (
	LOG   = true
	DEBUG = true
)

// Panics on the message.
func Panic(msg ...interface{}) {
	panic(fmt.Sprint(msg...))
}

// Panics if err is not nil
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Logs the error of non-nil.
func LogErr(err error) {
	if err != nil {
		Log("[error]", err)
	}
}

// Log message.
func Log(msg ...interface{}) {
	if LOG {
		log.Println(msg...)
	}
}

// Debug message.
func Debug(msg ...interface{}) {
	if DEBUG {
		Log(msg...)
	}
}

// Error message
func Error(msg ...interface{}) {
	Log(msg...)
}