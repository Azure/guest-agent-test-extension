package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

var (
	infoOperation    = "INFO"
	warningOperation = "WARNING"
	errorOperation   = "ERROR"
)

type customInfoLogger struct {
	infoLogger *log.Logger
}

func (iLogger *customInfoLogger) Print(v ...interface{}) {
	iLogger.infoLogger.Print(formatMessage(infoOperation, fmt.Sprint(v...)))
}
func (iLogger *customInfoLogger) Println(v ...interface{}) {
	iLogger.infoLogger.Println(formatMessage(infoOperation, fmt.Sprint(v...)))
}
func (iLogger *customInfoLogger) Printf(format string, v ...interface{}) {
	iLogger.infoLogger.Printf(formatMessage(infoOperation, fmt.Sprintf(format, v...)))
}

type customWarningLogger struct {
	warningLogger *log.Logger
}

func (wLogger *customWarningLogger) Print(v ...interface{}) {
	wLogger.warningLogger.Print(formatMessage(warningOperation, fmt.Sprint(v...)))
}
func (wLogger *customWarningLogger) Println(v ...interface{}) {
	wLogger.warningLogger.Println(formatMessage(warningOperation, fmt.Sprint(v...)))
}
func (wLogger *customWarningLogger) Printf(format string, v ...interface{}) {
	wLogger.warningLogger.Printf(formatMessage(warningOperation, fmt.Sprintf(format, v...)))
}

type customErrorLogger struct {
	errorLogger *log.Logger
}

func (eLogger *customErrorLogger) Print(v ...interface{}) {
	eLogger.errorLogger.Print(formatMessage(errorOperation, fmt.Sprint(v...)))
}
func (eLogger *customErrorLogger) Println(v ...interface{}) {
	eLogger.errorLogger.Println(formatMessage(errorOperation, fmt.Sprint(v...)))
}
func (eLogger *customErrorLogger) Printf(format string, v ...interface{}) {
	eLogger.errorLogger.Printf(formatMessage(errorOperation, fmt.Sprintf(format, v...)))
}

func formatMessage(messageType string, message string) string {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999999Z")
	_, filePath, line, _ := runtime.Caller(2)
	shortFileArr := strings.Split(filePath, "/")
	shortFile := shortFileArr[len(shortFileArr)-1]

	return fmt.Sprintf("[%s] [%s] [%s:%d] [%s]: %s", currentTime, version, shortFile, line, messageType, message)
}
