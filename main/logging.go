package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/Azure/azure-extension-foundation/settings"
	"github.com/pkg/errors"
)

type logPrefixType string

const (
	infoOperation    logPrefixType = "INFO"
	warningOperation logPrefixType = "WARNING"
	errorOperation   logPrefixType = "ERROR"
)

type loggerMode string

const (
	generalLoggerMode   loggerMode = "generalLoggerMode"
	operationLoggerMode loggerMode = "operationLoggerMode"
)

type customLogger struct {
	logger        *log.Logger
	messagePrefix logPrefixType
	mode          loggerMode
}

func (l *customLogger) Print(v ...interface{}) {
	l.logger.Print(formatLoggingMessage(l.mode, l.messagePrefix, fmt.Sprint(v...)))
}
func (l *customLogger) Println(v ...interface{}) {
	l.logger.Println(formatLoggingMessage(l.mode, l.messagePrefix, fmt.Sprint(v...)))
}
func (l *customLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(formatLoggingMessage(l.mode, l.messagePrefix, fmt.Sprintf(format, v...)))
}

func formatLoggingMessage(loggerType loggerMode, messagePrefix logPrefixType, message string) string {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	switch loggerType {
	case generalLoggerMode:
		//Caller(1) is the customlogger Print function in this file
		//Caller(2) is the line that called the logger's print function
		_, filePath, line, _ := runtime.Caller(2)

		//Gets the line number and caller of the log message to identify what line logged
		_, shortFile := path.Split(filePath)

		return fmt.Sprintf("[%s] [%s] [%s:%d] [%s]: %s", currentTime, version, shortFile, line, messagePrefix, message)
	case operationLoggerMode:
		// Message prefix for the operation logger is always INFO
		return fmt.Sprintf("[%s] [%s]: [Seq Num: %d] [%s]: %s", currentTime, version, environmentMrSeq, messagePrefix, message)
	default:
		fmt.Printf("Unable to determine logger type %s", loggerType)
		return fmt.Sprintf("[%s]: ERROR FORMATTING MESSAGE, UNKNOWN LOGGER TYPE: %s", currentTime, message)
	}
}

func initLoggingFilepath(logfileLogName string) (file *os.File, err error) {
	handlerEnv, handlerEnvErr := settings.GetHandlerEnvironment()
	var logfileFilepath string

	// Create the handler environment specified folder if it does not exist
	if handlerEnvErr != nil {
		logfileFilepath = logfileLogName
		fmt.Printf("Error opening handler environment %+v", handlerEnvErr)
	} else {
		if _, err := os.Stat(handlerEnv.HandlerEnvironment.LogFolder); os.IsNotExist(err) {
			os.Mkdir(handlerEnv.HandlerEnvironment.LogFolder, os.ModeDir)
		}
		logfileFilepath = path.Join(handlerEnv.HandlerEnvironment.LogFolder, logfileLogName)
	}

	// Create the log file if it does not exist, otherwise just open it
	file, err = os.OpenFile(logfileFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return file, err
}

func initGeneralLogging() (file *os.File, err error) {
	logfileLogName := extensionName + ".log"
	file, err = initLoggingFilepath(logfileLogName)

	//Sample: [2020-08-18T20:29:16.079Z] [1.0.0.0] [main.go:148] [INFO]: Test1
	infoLogger = customLogger{log.New(io.MultiWriter(file, os.Stdout), "", 0), infoOperation, generalLoggerMode}
	warningLogger = customLogger{log.New(io.MultiWriter(file, os.Stderr), "", 0), warningOperation, generalLoggerMode}
	errorLogger = customLogger{log.New(io.MultiWriter(file, os.Stderr), "", 0), errorOperation, generalLoggerMode}

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", logfileLogName)
	}
	return file, err
}

func initOperationLogging() (file *os.File, err error) {
	operationLogfileLogName := "operations-" + version + ".log"
	file, err = initLoggingFilepath(operationLogfileLogName)

	//Sample: [2020-09-05T00:23:40.578Z] [1.0.0.1]: [Seq Num: -1] [INFO]: install
	operationLogger = customLogger{log.New(file, "", 0), infoOperation, operationLoggerMode}

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", operationLogfileLogName)
	}
	return file, err
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
