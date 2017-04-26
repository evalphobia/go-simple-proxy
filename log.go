package main

import (
	"io/ioutil"
	"log"
	"os"
)

// logger list
var (
	infoLogger  = log.New(os.Stdout, "[PROXY INFO] ", log.LstdFlags)
	errLogger   = log.New(os.Stderr, "[PROXY ERROR] ", log.LstdFlags)
	debugLogger = log.New(ioutil.Discard, "[PROXY DEBUG] ", log.LstdFlags)
)

// disableLog disables logs.
func disableLog() {
	infoLogger.SetOutput(ioutil.Discard)
	debugLogger.SetOutput(ioutil.Discard)
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
