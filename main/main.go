package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/azure-extension-foundation/sequence"
	"github.com/Azure/azure-extension-foundation/settings"
	"github.com/Azure/azure-extension-foundation/status"
	"github.com/pkg/errors"
)

var (
	versionMajor    = "1"
	versionMinor    = "0"
	versionBuild    = "0"
	versionRevision = "2"
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
)

const (
	// Exit Codes
	successfulExecution   = iota // 0
	generalExitError             // 1
	commandNotFoundError         // 2
	logfileNotOpenedError        // 3
	mrSeqNotFoundError           // 4
	shouldNotRunError            // 5
	seqNumberSetError            // 6
	statusReportingError         // 7
	settingsNotFoundError        // 8
	versionMismatchError         // 9
	jsonParsingError             // 10
)

// extension specific PublicSettings
type PublicSettings struct {
	Script   string   `json:"script"`
	FileURLs []string `json:"fileUris"`
}

// extension specific ProtectedSettings
type ProtectedSettings struct {
	SecretString       string   `json:"secretString"`
	SecretScript       string   `json:"secretScript"`
	FileURLs           []string `json:"fileUris"`
	StorageAccountName string   `json:"storageAccountName"`
	StorageAccountKey  string   `json:"storageAccountKey"`
}

func install() {
	operation := "install"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Printf("[Seq Num: %d] [operation: %s]", environmentMrSeq, operation)

	err := status.ReportTransitioning(environmentMrSeq, operation, "installation in progress")
	infoLogger.Println("Installation in progress")
	if err != nil {
		errorLogger.Printf("Status reporting error: %+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "installation is complete")
	if err != nil {
		errorLogger.Printf("Status reporting error: %+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("Installation is complete")
	os.Exit(successfulExecution)
}

func enable() {
	operation := "enable"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Printf("[Seq Num: %d] [operation: %s]", environmentMrSeq, operation)

	err := status.ReportTransitioning(environmentMrSeq, operation, "enabling in progress")
	infoLogger.Println("enabling in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	var publicSettings PublicSettings
	var protectedSettings ProtectedSettings

	err = settings.GetExtensionSettings(environmentMrSeq, &publicSettings, &protectedSettings)
	if err != nil {
		status.ReportError(environmentMrSeq, operation, err.Error())
		errorLogger.Printf("%+v", err)
		os.Exit(settingsNotFoundError)
	}
	infoLogger.Printf("Public Settings: %v \t Protected Settings: %v", publicSettings, protectedSettings)

	err = status.ReportSuccess(environmentMrSeq, operation, "enabling is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("enabling is complete")
	os.Exit(successfulExecution)
}

func disable() {
	operation := "disable"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Printf("[Seq Num: %d] [operation: %s]", environmentMrSeq, operation)

	err := status.ReportTransitioning(environmentMrSeq, operation, "disabling in progress")
	infoLogger.Println("disabling in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "disabling is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("disabling is complete")
	os.Exit(successfulExecution)
}

func uninstall() {
	operation := "uninstall"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Printf("[Seq Num: %d] [operation: %s]", environmentMrSeq, operation)

	err := status.ReportTransitioning(environmentMrSeq, operation, "uninstallation in progress")
	infoLogger.Println("uninstallation in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "uninstallation is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("uninstallation is complete")
	os.Exit(successfulExecution)
}

func update() {
	operation := "update"
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)
	operationLogger.Printf("[Seq Num: %d] [operation: %s]", environmentMrSeq, operation)

	err := status.ReportTransitioning(environmentMrSeq, operation, "updating in progress")
	infoLogger.Println("updating in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "updating is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("updating is complete")
	os.Exit(successfulExecution)
}

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
	var jsonData map[string][]string
	json.Unmarshal([]byte(byteValue), &jsonData)

	fmt.Print(jsonData["failCommands"])
	return nil
}

func main() {
	generalFile, operationFile, generalErr, operationErr := initAllLogging()
	if generalErr != nil {
		fmt.Printf("Error opening general logfile %+v", generalErr)
		os.Exit(logfileNotOpenedError)
	}
	defer generalFile.Close()
	if operationErr != nil {
		fmt.Printf("Error opening operations logfile %+v", operationErr)
		os.Exit(logfileNotOpenedError)
	}
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
