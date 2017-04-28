package main

import (
	"io/ioutil"
	"log"
	"os"
)

// logger list
var (
	errLogger   = log.New(os.Stderr, "[PROXY ERROR] ", log.LstdFlags)
	infoLogger  = log.New(ioutil.Discard, "[PROXY INFO] ", log.LstdFlags)
	debugLogger = log.New(ioutil.Discard, "[PROXY DEBUG] ", log.LstdFlags)
)

// enableInfoLog enables info logs.
func enableInfoLog() {
	infoLogger.SetOutput(os.Stdout)
}

// enableDebugLog enables debug logs.
func enableDebugLog() {
	debugLogger.SetOutput(os.Stdout)
}

// loggingInfo outputs info log.
func loggingInfo(log string, v ...interface{}) {
	infoLogger.Printf(log+"\n", v...)
}

// loggingError outputs error log.
func loggingError(log string, v ...interface{}) {
	errLogger.Printf(log+"\n", v...)
}

// loggingDebug outputs debug log.
func loggingDebug(log string, v ...interface{}) {
	debugLogger.Printf(log+"\n", v...)
}
