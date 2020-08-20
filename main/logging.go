package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/Azure/azure-extension-foundation/settings"
	"github.com/pkg/errors"
)

var (
	infoOperation    = "INFO"
	warningOperation = "WARNING"
	errorOperation   = "ERROR"
)

type customGeneralLogger struct {
	logger    *log.Logger
	operation string
}

func (l *customGeneralLogger) Print(v ...interface{}) {
	l.logger.Print(formatLoggingMessage(l.operation, fmt.Sprint(v...)))
}

func (l *customGeneralLogger) Println(v ...interface{}) {
	l.logger.Println(formatLoggingMessage(l.operation, fmt.Sprint(v...)))
}

func (l *customGeneralLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(formatLoggingMessage(l.operation, fmt.Sprintf(format, v...)))
}

func formatLoggingMessage(messageType string, message string) string {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

	//Caller(1) is the customlogger Print function in this file
	//Caller(2) is the line that called the logger's print function
	_, filePath, line, _ := runtime.Caller(2)

	//Gets the line number and caller of the log message to identify what line logged
	shortFileArr := strings.Split(filePath, "/")
	shortFile := shortFileArr[len(shortFileArr)-1]

	return fmt.Sprintf("[%s] [%s] [%s:%d] [%s]: %s", currentTime, version, shortFile, line, messageType, message)
}

type customOperationLogger struct {
	logger *log.Logger
}

func (l *customOperationLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(formatOperationMessage(fmt.Sprintf(format, v...)))
}

func formatOperationMessage(message string) string {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	return fmt.Sprintf("[%s] [%s]: %s", currentTime, version, message)
}

func initAllLogging() (generalFile *os.File, operationFile *os.File, generalErr error, operationErr error) {
	generalFile, generalErr = initGeneralLogging()
	if generalErr != nil {
		fmt.Printf("Error opening the general logfile. %+v", generalErr)
		generalErr = errors.Wrap(generalErr, "Failed to open general logfile")
	}

	operationFile, operationErr = initOperationLogging()
	if operationErr != nil {
		fmt.Printf("Error opening the operation logfile. %+v", operationErr)
		operationErr = errors.Wrap(operationErr, "Failed to open operation logfile")
	}
	return generalFile, operationFile, generalErr, operationErr
}

func initGeneralLogging() (*os.File, error) {
	handlerEnv, handlerEnvErr := settings.GetHandlerEnvironment()

	logfileLogName := extensionName + "-" + version + ".log"
	if handlerEnvErr != nil {
		generalLogfile = logfileLogName
	} else {
		if _, err := os.Stat(handlerEnv.HandlerEnvironment.LogFolder); os.IsNotExist(err) {
			os.Mkdir(handlerEnv.HandlerEnvironment.LogFolder, os.ModeDir)
		}
		generalLogfile = path.Join(handlerEnv.HandlerEnvironment.LogFolder, logfileLogName)
	}

	file, err := os.OpenFile(generalLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", logfile)
	}

	//Sample: [2020-08-18T20:29:16.079902Z] [1.0.0.0] [main.go:148] [INFO]: Test1
	infoLogger = customGeneralLogger{log.New(io.MultiWriter(file, os.Stdout), "", 0), infoOperation}
	warningLogger = customGeneralLogger{log.New(io.MultiWriter(file, os.Stderr), "", 0), warningOperation}
	errorLogger = customGeneralLogger{log.New(io.MultiWriter(file, os.Stderr), "", 0), errorOperation}

	if handlerEnvErr != nil {
		errorLogger.Printf("Error opening handler environment %+v", handlerEnvErr)
	}
	return file, nil
}

func initOperationLogging() (*os.File, error) {
	handlerEnv, handlerEnvErr := settings.GetHandlerEnvironment()

	operationLogfileLogName := "operation-" + version + ".log"
	if handlerEnvErr != nil {
		operationLogfile = operationLogfileLogName
	} else {
		if _, err := os.Stat(handlerEnv.HandlerEnvironment.LogFolder); os.IsNotExist(err) {
			os.Mkdir(handlerEnv.HandlerEnvironment.LogFolder, os.ModeDir)
		}
		operationLogfile = path.Join(handlerEnv.HandlerEnvironment.LogFolder, operationLogfileLogName)
	}

	file, err := os.OpenFile(operationLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", operationLogfile)
	}

	//Sample: [2020-08-18T20:29:16.079902Z] [1.0.0.0] [main.go:148] [INFO]: Test1
	operationLogger = customOperationLogger{log.New(io.MultiWriter(file, os.Stdout), "", 0)}

	if handlerEnvErr != nil {
		errorLogger.Printf("Error opening handler environment %+v", handlerEnvErr)
	}
	return file, nil
}
