package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-extension-foundation/sequence"
	"github.com/Azure/azure-extension-foundation/settings"
	"github.com/Azure/azure-extension-foundation/status"
	"github.com/pkg/errors"
)

var (
	versionMajor    = "1"
	versionMinor    = "0"
	versionBuild    = "0"
	versionRevision = "1"
	version         = fmt.Sprintf("%s.%s.%s.%s", versionMajor, versionMinor, versionBuild, versionRevision)

	extensionMrSeq   int
	environmentMrSeq int

	generalLogfile   string
	operationLogfile string

	extensionName = "GuestAgentTestExtension"

	extensionConfiguration                                  extensionConfigurationStruct
	failCommands                                            []failCommandsStruct
	infoLogger, warningLogger, errorLogger, operationLogger customLogger

	// Execution errors that are encountered during execution are stored and then reported at the end
	executionErrors []string
	// Any exit code specified like in failCommand that should be used to exit
	intendedExitCode = successfulExecution
)

const (
	// Pre-determined exeit codes
	successfulExecution   = iota // 0
	generalExitError             // 1
	commandNotFoundError         // 2
	logfileNotOpenedError        // 3
)

// TODO future runtime configuration can be added in this struct
type extensionConfigurationStruct struct {
	FailCommands []failCommandsStruct `json:"failCommands"`
}

// Format of the failCommands in the runtime configuration json file
type failCommandsStruct struct {
	Command               string `json:"command"`
	ErrorMessage          string `json:"errorMessage"`
	ExitCode              string `json:"exitCode"`
	ReportStatusCorrectly string `json:"reportStatusCorrectly"`
}

// extension specific PublicSettings
type extensionPublicSettings struct {
	Name string `json:"name"`
}

// extension specific ProtectedSettings
type extensionPrivateSettings struct {
	SecretString string `json:"secretString"`
}

// This is an implementation from the golang extension library, but these consts are not exported so are reproduced here
type extensionStatus string

const (
	statusTransitioning extensionStatus = "transitioning"
	statusError         extensionStatus = "error"
	statusSuccess       extensionStatus = "success"
)

// Reports the status with error handling for an operation
func reportStatus(statusType extensionStatus, operation string, message string) {
	var err error
	isFailCommand := 0

	// TODO: This could probably be refactored but unfortunately I do not have time to finish the implementation
	// if failCommands were instead stored as a map[string] []string where the accessing string was operation,
	// we would not have to cycle through all the fail commands every time and could just read the map directly
	for _, failCommand := range failCommands {
		if failCommand.Command == operation && statusType == statusSuccess {
			errorLogger.Printf("%s failed with message: %s Expected exitCode: %s", failCommand.Command, failCommand.ErrorMessage, failCommand.ExitCode)

			errorLogger.Printf("%s will report status correctly before exiting: %s", failCommand.Command, failCommand.ReportStatusCorrectly)

			// Only report an error in status if the failCommand should report status correctly
			if failCommand.ReportStatusCorrectly == "true" {
				err := status.ReportError(environmentMrSeq, operation, failCommand.ErrorMessage)
				if err != nil {
					errorMessage := fmt.Sprintf("Status reporting error: %+v", err)
					errorLogger.Println(errorMessage)
					executionErrors = append(executionErrors, errorMessage)
				}
			}

			// Get the exit code and save it for when os.Exit will be called
			if exitCode, err := strconv.Atoi(failCommand.ExitCode); err == nil {
				intendedExitCode = exitCode
				isFailCommand = 1
			} else {
				errorMessage := fmt.Sprintf("Unable to use provided exit code %+v", err)
				errorLogger.Println(errorMessage)
				executionErrors = append(executionErrors, errorMessage)
			}
		}
	}

	switch statusType {
	case statusSuccess:
		// no success status for fail commands, status will have already been updated
		// TODO this implementaiton is only correct for a failure after transitioning has started.
		// if we want to support a scenario where the extension somehow failed to start transitioning
		// this would have to be changed
		if isFailCommand == 0 {
			err = status.ReportSuccess(environmentMrSeq, operation, message)
			infoLogger.Println(message)
		}
	case statusTransitioning:
		err = status.ReportTransitioning(environmentMrSeq, operation, message)
		infoLogger.Println(message)
	case statusError:
		err = status.ReportError(environmentMrSeq, operation, message)
		errorLogger.Println(message)
	default:
		warningLogger.Println("Status report type not recognized")
	}

	if err != nil {
		errorMessage := fmt.Sprintf("Status reporting error: %+v", err)
		errorLogger.Println(errorMessage)
		executionErrors = append(executionErrors, errorMessage)
	}
}

// Basic functionality for commands that do not have much special behavior (like enable)
func testCommand(operation string) {
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Println(operation)

	reportStatus(statusTransitioning, operation, fmt.Sprintf("%s in progress", operation))
	reportStatus(statusSuccess, operation, fmt.Sprintf("%s completed successfully", operation))
}

// Handles how the extension exits when execution is done. Prints out errors if there are any and uses
// any exit codes that were specified
func reportExecutionStatus() {
	if executionErrors == nil {
		infoLogger.Printf("Exiting with Code: %d", intendedExitCode)
		os.Exit(intendedExitCode)
	} else {
		errorMessage := strings.Join(executionErrors, "\n")
		errorLogger.Println(errorMessage)
		os.Exit(generalExitError)
	}
}

func install() {
	operation := "install"
	testCommand(operation)
}

// Enable prints out the name provided in the public settings
func enable() {
	operation := "enable"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Println(operation)
	reportStatus(statusTransitioning, operation, fmt.Sprintf("%s in progress", operation))

	var publicSettings extensionPublicSettings
	var protectedSettings extensionPrivateSettings

	err := settings.GetExtensionSettings(environmentMrSeq, &publicSettings, &protectedSettings)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting settings: %+v", err)
		errorLogger.Println(errorMessage)
		executionErrors = append(executionErrors, errorMessage)
		reportStatus(statusError, operation, fmt.Sprintf("%s failed due to inability to getting settings", operation))
		return
	}
	infoLogger.Printf("Public Settings: %v \t Protected Settings: %v", publicSettings, protectedSettings)
	infoLogger.Printf("Provided Name is: %s", publicSettings.Name)

	reportStatus(statusSuccess, operation, fmt.Sprintf("%s completed successfully", operation))
}

func disable() {
	operation := "disable"
	testCommand(operation)
}

func uninstall() {
	operation := "uninstall"
	testCommand(operation)
}

func update() {
	operation := "update"
	testCommand(operation)
}

// Parses the JSON file using the predetermined structs for formatting
func parseJSON(filename string) error {
	//	Open the provided file
	jsonFile, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to open \"%s\"", filename)
	}
	infoLogger.Println("File opened successfully")

	// Defer file closing until parseJSON() returns
	defer jsonFile.Close()

	//	Unmarshall the bytes from the JSON file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &extensionConfiguration)

	failCommands = extensionConfiguration.FailCommands

	// Future configuation paramters should be set here

	return nil
}

func main() {
	generalFile, operationFile, loggingErr := initAllLogging()
	if loggingErr != nil {
		fmt.Printf("Error opening general logfile %+v", loggingErr)
		os.Exit(logfileNotOpenedError)
	}
	defer generalFile.Close()
	defer operationFile.Close()
	// TODO: The file won't open if init logging throws an error, but file.close can also
	// have errors related to disk writing delays. This works well enough, but since this is
	// a testing extension, it might be worth adding additional error handling

	envExtensionVersion := os.Getenv("AZURE_GUEST_AGENT_EXTENSION_VERSION")
	if envExtensionVersion != "" && envExtensionVersion != version {
		warningLogger.Printf("Internal version %s does not match with environment variable version %s",
			version, envExtensionVersion)
	}
	var err error

	extensionMrSeq, environmentMrSeq, err = sequence.GetMostRecentSequenceNumber()
	if err != nil {
		warningLogger.Printf("Error getting sequence number %+v", err)
		extensionMrSeq = -1
		environmentMrSeq = -1
	}
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)

	// Command line flags that are currently supported
	commandStringPtr := flag.String("command", "", "Valid commands are install, enable, update, disable and uninstall. Usage: --command=install")
	parseJSONPtr := flag.String("jsonfile", "", "Path to the JSON file loction. Usage --jsonfile=\"test.json\"")

	// Trigger parsing of the command flag and then run the corresponding command
	flag.Parse()

	err = parseJSON(*parseJSONPtr)
	if err != nil {
		errorLogger.Printf("Error parsing provided JSON file: %+v", err)
	}

	switch *commandStringPtr {
	case "disable":
		disable()
	case "uninstall":
		uninstall()
	case "install":
		install()
	case "enable":
		enable()
	case "update":
		update()
	case "":
		warningMessage := "No --command flag provided"
		warningLogger.Println(warningMessage)
	default:
		warningMessage := fmt.Sprintf("Command \"%s\" not recognized", *commandStringPtr)
		warningLogger.Println(warningMessage)
	}

	reportExecutionStatus()
}
