package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-extension-foundation/sequence"
	"github.com/Azure/azure-extension-foundation/settings"
	"github.com/Azure/azure-extension-foundation/status"
	"github.com/pkg/errors"
)

var (
	versionMajor    = "1"
	versionMinor    = "0"
	versionBuild    = "1"
	versionRevision = "0"
	version         = fmt.Sprintf("%s.%s.%s.%s", versionMajor, versionMinor, versionBuild, versionRevision)

	extensionMrSeq   int
	environmentMrSeq int

	// Logging is currently set up to create/add to the logile in the directory from where the binary is executed
	// TODO Read this in from Handler Env
	logfile        string
	logfileLogName = "GuestAgentTestExtension-" + version + ".log"

	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
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
	infoLogger.Printf("Sequence number is %d", environmentMrSeq)

	err := status.ReportTransitioning(environmentMrSeq, operation, "installation in progress")
	if err != nil {
		errorLogger.Printf("Installation error: %+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("Installation in progress")

	err = status.ReportSuccess(environmentMrSeq, operation, "installation is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}
	infoLogger.Println("Installation is complete")
	os.Exit(successfulExecution)
}

func enable() {
	operation := "enable"

	err := status.ReportTransitioning(environmentMrSeq, operation, "enabling in progress")
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

	os.Exit(successfulExecution)
}

func disable() {
	operation := "disable"

	err := status.ReportTransitioning(environmentMrSeq, operation, "disabling in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "disabling is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	os.Exit(successfulExecution)
}

func uninstall() {
	operation := "uninstall"

	err := status.ReportTransitioning(environmentMrSeq, operation, "uninstallation in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "uninstallation is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	os.Exit(successfulExecution)
}

func update() {
	operation := "update"

	err := status.ReportTransitioning(environmentMrSeq, operation, "updating in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

	err = status.ReportSuccess(environmentMrSeq, operation, "updating is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(statusReportingError)
	}

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
	// 	TODO: If we know the exact format, we can read the JSON into a struct which might be cleaner
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	// Get the map with the name "keys"
	keys := result["keys"].(map[string]interface{})

	//	Parse each key value and reverse the string by appending characters backwards
	for key, value := range keys {
		reverseValue := ""
		for _, val := range value.(string) {
			reverseValue = string(val) + reverseValue
		}

		infoLogger.Println(key, reverseValue)
	}
	return nil
}

/* 	Open the logfile and configure the loggers that will be used
*
*	The main difference between types of loggers is the label (eg INFO) and additional data provided .
 */
func initLogging() (*os.File, error) {
	extensionPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	externalVersion := strings.Split(strings.Split(filepath.Base(extensionPath), "-")[1], ".")
	// If the revision number (4th digit) is not present, add a 0 to the end
	if len(externalVersion) < 4 {
		externalVersion = append(externalVersion, "0")
	}
	externalVersionString := strings.Join(externalVersion, ".")

	if externalVersionString != version {
		fmt.Printf("Current internal extension version %s does not match directory version %s using directory version of extension\n", version, externalVersionString)
		version = externalVersionString
	}

	handlerEnv, err := settings.GetHandlerEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open handler environment")
	}

	logfile = path.Join(handlerEnv.HandlerEnvironment.LogFolder, logfileLogName)

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create/open %s", logfile)
	}

	// Log UTC Time, Date, Time, (w/ microseconds), line number, and move message prefix
	loggerFlags := log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix

	// Sample out: 2020/07/31 21:47:19.153535 main.go:145: Version: 1.0.0.0 INFO: Hello World!
	infoLogger = log.New(io.MultiWriter(file, os.Stdout), "Version: "+version+" INFO: ", loggerFlags)

	// Sample out: 2020/07/31 21:47:19.153535 main.go:145: Version: 1.0.0.0 WARNING: Hello World!
	warningLogger = log.New(io.MultiWriter(file, os.Stderr), "Version: "+version+" WARNING: ", loggerFlags)

	// Sample out: 2020/07/31 21:47:19.153535 main.go:145: Version: 1.0.0.0 ERROR: Hello World!
	errorLogger = log.New(io.MultiWriter(file, os.Stderr), "Version: "+version+" ERROR: ", loggerFlags)

	return file, nil
}

func main() {
	file, err := initLogging()
	if err != nil {
		fmt.Printf("Error opening the provided logfile. %+v", err)
		os.Exit(logfileNotOpenedError)
	}
	//TODO: The file won't open if init logging throws an error, but file.close can also
	//have errors related to disk writing delays. Will update with more robust error handling
	//but for now this works well enough
	defer file.Close()

	extensionMrSeq, environmentMrSeq, err = sequence.GetMostRecentSequenceNumber()
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(mrSeqNotFoundError)
	}
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrSeq, environmentMrSeq)

	shouldRun := sequence.ShouldBeProcessed(extensionMrSeq, environmentMrSeq)
	if !shouldRun {
		errorLogger.Printf("environment mrseq has already been processed by extension (environment mrseq : %v, extension mrseq : %v)\n", environmentMrSeq, extensionMrSeq)
		os.Exit(shouldNotRunError)
	}
	infoLogger.Printf("Extension should run: %t", shouldRun)

	err = sequence.SetExtensionMostRecentSequenceNumber(environmentMrSeq)
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(seqNumberSetError)
	}

	// Command line flags that are currently supported
	commandStringPtr := flag.String("command", "", "Valid commands are install, enable, update, disable and uninstall. Usage: --command=install")
	parseJSONPtr := flag.String("jsonfile", "", "Path to the JSON file loction. Usage --jsonfile=\"test.json\"")

	// Trigger parsing of the command flag and then run the corresponding command
	flag.Parse()
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
	default:
		warningLogger.Printf("Command \"%s\" not recognized", *commandStringPtr)
		os.Exit(commandNotFoundError)
	}
	//TODO this is being deprecated since it was proof of concept and the method will be repurposed
	// Parse the provided JSON file if there is one
	if *parseJSONPtr != "" {
		err := parseJSON(*parseJSONPtr)
		if err != nil {
			// Gives full traceback: Sample:
			/* 2020/07/31 22:28:54.962003 main.go:144: Version: 1.0.0.0 ERROR: open tet.json: The system cannot find the file specified.
			Failed to open "tet.json"
			main.parseJSON
				C:/Users/t-etfali/Documents/GuestAgentExtension/guest-agent-test-extension/main.go:52
			main.main
				C:/Users/t-etfali/Documents/GuestAgentExtension/guest-agent-test-extension/main.go:141
			runtime.main
				c:/go/src/runtime/proc.go:203
			runtime.goexit
				c:/go/src/runtime/asm_amd64.s:1373
			*/

			errorLogger.Printf("%+v", err)
			os.Exit(generalExitError)
		}
	}
}
