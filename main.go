package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/Azure/azure-extension-foundation/sequence"
	"github.com/Azure/azure-extension-foundation/settings"
	"github.com/Azure/azure-extension-foundation/status"
	"github.com/pkg/errors"
)

var (
	version            = "1.0.0.0"
	extensionShortName = "GATestExt"

	// Logging is currently set up to create/add to the logile in the directory from where the binary is executed
	// TODO Read this in from Handler Env
	logfile        string
	logfileLogName = "guest-agent-test-extension.log"

	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

func install() {
	infoLogger.Println("Installed Succesfully.")
}

func enable() {
	infoLogger.Println("Enabled Successfully.")
}

func disable() {
	infoLogger.Println("Disabled Successfully.")
}

func uninstall() {
	infoLogger.Println("Uninstalled Successfully.")
}

func update() {
	infoLogger.Println("Updated Successfully.")
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
	handlerEnv, err := settings.GetHandlerEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open handler environment")
	}

	logfile = handlerEnv.HandlerEnvironment.LogFolder + "/" + logfileLogName

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create/open specified log file")
	}

	/* Log UTC Time, Date, Time, (w/ microseconds), line number, and make message prefix come right
	before the message
	*/
	loggerFlags := log.LUTC | log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix

	// Sample out: 2020/07/31 21:47:19.153535 main.go:145: Version: 1.0.0.0 INFO: Hello World!
	infoLogger = log.New(io.MultiWriter(file, os.Stdout), "Version: "+version+" INFO: ", loggerFlags)

	// Sample out: 2020/07/31 21:47:19.153535 main.go:145: Version: 1.0.0.0 WARNING: Hello World!
	warningLogger = log.New(io.MultiWriter(file, os.Stderr), "Version: "+version+" WARNING: ", loggerFlags)

	// Sample out: 2020/07/31 21:47:19.153535 main.go:145: Version: 1.0.0.0 ERROR: Hello World!
	errorLogger = log.New(io.MultiWriter(file, os.Stderr), "Version: "+version+" ERROR: ", loggerFlags)

	return file, nil
}

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

func main() {
	file, err := initLogging()
	if err != nil {
		fmt.Println("Error opening the provided logfile.")
		os.Exit(-1)
	}
	defer file.Close()

	extensionMrseq, environmentMrseq, err := sequence.GetMostRecentSequenceNumber()
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(-1)
	}
	infoLogger.Printf("Extension MrSeq: %d, Environment MrSeq: %d", extensionMrseq, environmentMrseq)

	shouldRun := sequence.ShouldBeProcessed(extensionMrseq, environmentMrseq)
	if !shouldRun {
		errorLogger.Printf("environment mrseq has already been processed by extension (environment mrseq : %v, extension mrseq : %v)\n", environmentMrseq, extensionMrseq)
		os.Exit(-1)
	}
	infoLogger.Printf("Extension should run: %t", shouldRun)

	err = sequence.SetExtensionMostRecentSequenceNumber(environmentMrseq)
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(-1)
	}

	err = status.ReportTransitioning(environmentMrseq, "install", "installation in progress")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(-1)
	}

	var publicSettings PublicSettings
	var protectedSettings ProtectedSettings
	err = settings.GetExtensionSettings(environmentMrseq, &publicSettings, &protectedSettings)
	if err != nil {
		status.ReportError(environmentMrseq, "install", err.Error())
		errorLogger.Printf("%+v", err)
		os.Exit(-1)
	}
	infoLogger.Printf("Public Settings: %v \t Protected Settings: %v", publicSettings, protectedSettings)

	err = status.ReportSuccess(environmentMrseq, "install", "installation is complete")
	if err != nil {
		errorLogger.Printf("%+v", err)
		os.Exit(-1)
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
	}

	// Parse the provided JSON file if there is one
	if *parseJSONPtr != "" {
		err := parseJSON(*parseJSONPtr)
		// TODO Add more robust error handling, Github-pkg-errors seems like a good candidate
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
			os.Exit(-1)
		}
	}
}
