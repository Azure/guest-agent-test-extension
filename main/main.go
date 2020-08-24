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
	versionRevision = "3"
	version         = fmt.Sprintf("%s.%s.%s.%s", versionMajor, versionMinor, versionBuild, versionRevision)

	extensionMrSeq   int
	environmentMrSeq int

	// Logging is currently set up to create/add to the logile in the directory from where the binary is executed
	// TODO Read this in from Handler Env
	generalLogfile   string
	operationLogfile string

	extensionName = "GuestAgentTestExtension"

	infoLogger, warningLogger, errorLogger customGeneralLogger
	operationLogger                        customOperationLogger

	executionErrors []string

	extensionConfiguration extensionConfigurationStruct
	failCommands           []failCommandsStruct
)

const (
	// Exit Codes
	successfulExecution   = iota // 0
	generalExitError             // 1
	commandNotFoundError         // 2
	logfileNotOpenedError        // 3
)

type extensionConfigurationStruct struct {
	FailCommands []failCommandsStruct `json:"failCommands"`
}

type failCommandsStruct struct {
	Command      string `json:"command"`
	ErrorMessage string `json:"errorMessage"`
	ExitCode     string `json:"exitCode"`
}

// extension specific PublicSettings
type extensionPublicSettings struct {
	Name string `json:"name"`
}

// extension specific ProtectedSettings
type extensionPrivateSettings struct {
	SecretString string `json:"secretString"`
}

func reportStatus(statusType string, operation string, message string) {
	switch statusType {
	case "success":
		err := status.ReportSuccess(environmentMrSeq, operation, message)
		infoLogger.Println(message)
		if err != nil {
			errorMessage := fmt.Sprintf("Status reporting error: %+v", err)
			errorLogger.Println(errorMessage)
			executionErrors = append(executionErrors, errorMessage)
		}
	case "transitioning":
		err := status.ReportTransitioning(environmentMrSeq, operation, message)
		infoLogger.Println(message)
		if err != nil {
			errorMessage := fmt.Sprintf("Status reporting error: %+v", err)
			errorLogger.Println(errorMessage)
			executionErrors = append(executionErrors, errorMessage)
		}
	case "error":
		err := status.ReportError(environmentMrSeq, operation, message)
		errorLogger.Println(message)
		if err != nil {
			errorMessage := fmt.Sprintf("Status reporting error: %+v", err)
			errorLogger.Println(errorMessage)
			executionErrors = append(executionErrors, errorMessage)
		}
	default:
		warningLogger.Println("Status report type not recognized")
	}

}

func checkForFailCommand(operation string) {
	for _, failCommand := range failCommands {
		if failCommand.Command == operation {
			reportStatus("error", operation, failCommand.ErrorMessage)
			if failCommand.ExitCode == "" {
				errorLogger.Printf("%s failed with message: %s, but will not exit since provided exitcode is %s", failCommand.Command, failCommand.ErrorMessage, failCommand.ExitCode)
			} else if exitCode, err := strconv.Atoi(failCommand.ExitCode); err == nil {
				errorLogger.Printf("%s failed with message: %s exitCode: %s", failCommand.Command, failCommand.ErrorMessage, failCommand.ExitCode)
				os.Exit(exitCode)
			} else {
				errorLogger.Printf("Unable to use provided exit code %+v", err)
				os.Exit(generalExitError)
			}

		}
	}
}

func testCommand(operation string) {
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Println(operation)
	reportStatus("transitioning", operation, fmt.Sprintf("%s in progress", operation))
	checkForFailCommand(operation)
	reportStatus("success", operation, fmt.Sprintf("%s completed successfully", operation))
}

func parseJSON(filename string) error {
	//	Open the provided file
	jsonFile, err := os.Open(filename)
	if err != nil {
		return errors.Wrapf(err, "Failed to open \"%s\"", filename)
	}
	infoLogger.Println("JSON File opened successfully")

	// Defer file closing until parseJSON() returns
	defer jsonFile.Close()

	//	Unmarshall the bytes from the JSON file
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &extensionConfiguration)

	failCommands = extensionConfiguration.FailCommands
	return nil
}

func reportExecutionStatus() {
	if executionErrors == nil {
		os.Exit(successfulExecution)
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

func enable() {
	operation := "enable"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Println(operation)
	reportStatus("transitioning", operation, fmt.Sprintf("%s in progress", operation))

	checkForFailCommand(operation)
	var publicSettings extensionPublicSettings
	var protectedSettings extensionPrivateSettings

	err := settings.GetExtensionSettings(environmentMrSeq, &publicSettings, &protectedSettings)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting settings: %+v", err)
		errorLogger.Println(errorMessage)
		executionErrors = append(executionErrors, errorMessage)
	}
	infoLogger.Printf("Public Settings: %v \t Protected Settings: %v", publicSettings, protectedSettings)
	infoLogger.Printf("Provided Name is: %s", publicSettings.Name)

	reportStatus("success", operation, fmt.Sprintf("%s completed successfully", operation))
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

func main() {
	generalFile, operationFile, loggingErr := initAllLogging()
	if loggingErr != nil {
		fmt.Printf("Error opening logfile %+v", loggingErr)
		os.Exit(logfileNotOpenedError)
	}
	defer generalFile.Close()
	defer operationFile.Close()
	//TODO: The file won't open if init logging throws an error, but file.close can also
	//have errors related to disk writing delays. Will update with more robust error handling
	//but for now this works well enough

	envExtensionVersion := os.Getenv("AZURE_GUEST_AGENT_EXTENSION_VERSION")
	if envExtensionVersion != "" && envExtensionVersion != version {
		warningLogger.Printf("Internal version %s does not match with environment variable version %s",
			version, envExtensionVersion)
	}

	extensionMrSeq, environmentMrSeq, err := sequence.GetMostRecentSequenceNumber()
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
		warningLogger.Println("No --command flag provided")
		os.Exit(commandNotFoundError)
	default:
		warningLogger.Printf("Command \"%s\" not recognized", *commandStringPtr)
		os.Exit(commandNotFoundError)
	}
}
