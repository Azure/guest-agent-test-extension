package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

var (
	version            = "1.0.0.0"
	extensionShortName = "GATestExt"

	// Logging is currently set up to create/add to the logile in the directory from where the binary is executed
	// TODO Read this in from Handler Env
	logfile = "guest-agent-test-extension.log"

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
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
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

func main() {
	// TODO Add more robust error handling
	file, err := initLogging()
	if err != nil {
		fmt.Println("Error opening the provided logfile.")
		os.Exit(1)
	}
	defer file.Close()

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
			// Gives full traceback
			errorLogger.Printf("%+v", err)
			os.Exit(1)
		}
	}
}
