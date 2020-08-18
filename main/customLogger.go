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

type customLogger struct {
	logger    *log.Logger
	operation string
}

func (l *customLogger) initLogger(loggerType string) {
	l.operation = loggerType
}

func (l *customLogger) Print(v ...interface{}) {
	l.logger.Print(formatMessage(l.operation, fmt.Sprint(v...)))
}

func (l *customLogger) Println(v ...interface{}) {
	l.logger.Println(formatMessage(l.operation, fmt.Sprint(v...)))
}

func (l *customLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(formatMessage(l.operation, fmt.Sprintf(format, v...)))
}

func formatMessage(messageType string, message string) string {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	_, filePath, line, _ := runtime.Caller(2)
	shortFileArr := strings.Split(filePath, "/")
	shortFile := shortFileArr[len(shortFileArr)-1]

	return fmt.Sprintf("[%s] [%s] [%s:%d] [%s]: %s", currentTime, version, shortFile, line, messageType, message)
}
