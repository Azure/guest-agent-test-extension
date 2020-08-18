package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

type customInfoLogger struct {
	infoLogger *log.Logger
}

func (iLogger *customInfoLogger) Write(bytes []byte) (int, error) {
	return iLogger.Write([]byte(formatMessage("INFO", string(bytes))))
}

func (iLogger *customInfoLogger) Print(v ...interface{}) {
	iLogger.infoLogger.Print(v...)
}
func (iLogger *customInfoLogger) Println(v ...interface{}) {
	iLogger.infoLogger.Println(v...)
}
func (iLogger *customInfoLogger) Printf(format string, v ...interface{}) {
	iLogger.infoLogger.Printf(format, v...)
}

type customWarningLogger struct {
	warningLogger *log.Logger
}

func (wLogger *customWarningLogger) Write(bytes []byte) (int, error) {
	return wLogger.Write([]byte(formatMessage("WARNING", string(bytes))))
}

type customErrorLogger struct {
	errorLogger *log.Logger
}

func (eLogger *customErrorLogger) Write(bytes []byte) (int, error) {
	return eLogger.Write([]byte(formatMessage("ERROR", string(bytes))))
}

func formatMessage(messageType string, message string) string {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999999Z")
	_, filePath, line, _ := runtime.Caller(1)
	shortFileArr := strings.Split(filePath, "/")
	shortFile := shortFileArr[len(shortFileArr)-1]

	return fmt.Sprintf("%s %s %s:%d [%s]: %s", currentTime, version, shortFile, line, messageType, message)

}
