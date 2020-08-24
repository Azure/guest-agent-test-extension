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

type customLogger struct {
	logger    *log.Logger
	operation string
}

type customGeneralLogger struct {
	customLogger
}

type customOperationLogger struct {
	customLogger
}

func (l *customGeneralLogger) Print(v ...interface{}) {
	l.logger.Print(formatLoggingMessage(*l, l.operation, fmt.Sprint(v...)))
}
func (l *customGeneralLogger) Println(v ...interface{}) {
	l.logger.Println(formatLoggingMessage(*l, l.operation, fmt.Sprint(v...)))
}
func (l *customGeneralLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(formatLoggingMessage(*l, l.operation, fmt.Sprintf(format, v...)))
}

func (l *customOperationLogger) Print(v ...interface{}) {
	l.logger.Print(formatLoggingMessage(*l, l.operation, fmt.Sprint(v...)))
}
func (l *customOperationLogger) Println(v ...interface{}) {
	l.logger.Println(formatLoggingMessage(*l, l.operation, fmt.Sprint(v...)))
}
func (l *customOperationLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(formatLoggingMessage(*l, l.operation, fmt.Sprintf(format, v...)))
}

func formatLoggingMessage(l interface{}, messageType string, message string) string {
	switch loggerType := l.(type) {
	case customGeneralLogger:
		currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")

		//Caller(1) is the customlogger Print function in this file
		//Caller(2) is the line that called the logger's print function
		_, filePath, line, _ := runtime.Caller(2)

		//Gets the line number and caller of the log message to identify what line logged
		shortFileArr := strings.Split(filePath, "/")
		shortFile := shortFileArr[len(shortFileArr)-1]

		return fmt.Sprintf("[%s] [%s] [%s:%d] [%s]: %s", currentTime, version, shortFile, line, messageType, message)

	case customOperationLogger:
		currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		return fmt.Sprintf("[%s] [%s]: [Seq Num: %d] [operation: %s]", currentTime, version, environmentMrSeq, message)

	default:
		fmt.Printf("Unable to determine logger type %s, string not formatted", loggerType)
		return message
	}
}

func initGeneralLogging() (file *os.File, err error) {
	handlerEnv, handlerEnvErr := settings.GetHandlerEnvironment()

	logfileLogName := extensionName + ".log"
	if handlerEnvErr != nil {
		generalLogfile = logfileLogName
	} else {
		if _, err := os.Stat(handlerEnv.HandlerEnvironment.LogFolder); os.IsNotExist(err) {
			os.Mkdir(handlerEnv.HandlerEnvironment.LogFolder, os.ModeDir)
		}
		generalLogfile = path.Join(handlerEnv.HandlerEnvironment.LogFolder, logfileLogName)
	}

	file, err = os.OpenFile(generalLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", generalLogfile)
	}

	//Sample: [2020-08-18T20:29:16.079902Z] [1.0.0.0] [main.go:148] [INFO]: Test1
	infoLogger = customGeneralLogger{customLogger{log.New(io.MultiWriter(file, os.Stdout), "", 0), infoOperation}}
	warningLogger = customGeneralLogger{customLogger{log.New(io.MultiWriter(file, os.Stderr), "", 0), warningOperation}}
	errorLogger = customGeneralLogger{customLogger{log.New(io.MultiWriter(file, os.Stderr), "", 0), errorOperation}}

	if handlerEnvErr != nil {
		errorLogger.Printf("Error opening handler environment %+v", handlerEnvErr)
	}
	return file, nil
}

func initOperationLogging() (file *os.File, err error) {
	handlerEnv, handlerEnvErr := settings.GetHandlerEnvironment()

	operationLogfileLogName := "operations-" + version + ".log"
	if handlerEnvErr != nil {
		operationLogfile = operationLogfileLogName
	} else {
		if _, err := os.Stat(handlerEnv.HandlerEnvironment.LogFolder); os.IsNotExist(err) {
			os.Mkdir(handlerEnv.HandlerEnvironment.LogFolder, os.ModeDir)
		}
		operationLogfile = path.Join(handlerEnv.HandlerEnvironment.LogFolder, operationLogfileLogName)
	}

	file, err = os.OpenFile(operationLogfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", operationLogfile)
	}

	//Sample: [2020-08-20T23:29:30.676Z] [1.0.0.2]: [Seq Num: 0] [operation: install]
	operationLogger = customOperationLogger{customLogger{log.New(file, "", 0), ""}}

	if handlerEnvErr != nil {
		errorLogger.Printf("Error opening handler environment %+v", handlerEnvErr)
	}
	return file, nil
}

func initAllLogging() (generalFile *os.File, operationFile *os.File, loggingErr error) {
	generalFile, loggingErr = initGeneralLogging()
	if loggingErr != nil {
		fmt.Printf("Error opening the general logfile. %+v", loggingErr)
		loggingErr = errors.Wrap(loggingErr, "Failed to open general logfile")
		return nil, nil, loggingErr
	}

	operationFile, loggingErr = initOperationLogging()
	if loggingErr != nil {
		fmt.Printf("Error opening the operation logfile. %+v", loggingErr)
		loggingErr = errors.Wrap(loggingErr, "Failed to open operation logfile")
		return nil, nil, loggingErr
	}

	return generalFile, operationFile, loggingErr
}
