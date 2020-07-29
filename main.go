package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	version            = "1.0.0.0"
	extensionShortName = "GATestExt"
	// Logging is currently set up to create/add to the logile in the directory from where the binary
	// is exectued. TODO: Add absolute filepath
	logfile     = "guest-agent-test-extension.log"
	logWriter   *os.File
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func install() {
	logInfo("Installed Succesfully.")
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

func parseJSON(s string) {
	//	Open the provided file
	jsonFile, err := os.Open(s)
	if err != nil {
		errorLogger.Println("File Not Found. JSON Parsing not performed.")
		os.Exit(1)
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

}

/* 	Open the logfile and configure the loggers that will be used
*
*	The main difference between types of loggers is the label (eg INFO) and additional data provided .
 */
func initLogging() {
	// Open the file as read only
	var err error
	logWriter, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		errorLogger.Println("Error opening file: ", logfile)
		log.Fatal(err)
	}
}

func logInfo(message string) {
	logMesssage := fmt.Sprintf(logfile, "%s %s %s INFO %s",
		time.Now().UTC().Format("2006-01-02T15:04:05.999999Z"), extensionShortName, version, message)

	fmt.Fprintf(logWriter, "%s\n", logMesssage)
}

func logWarning(message string) {
	logMesssage := fmt.Sprintf(logfile, "%s %s %s WARNING %s",
		time.Now().UTC().Format("2006-01-02T15:04:05.999999Z"), extensionShortName, version, message)

	fmt.Fprintf(logWriter, "%s\n", logMesssage)
}

func logError(message string) {
	logMesssage := fmt.Sprintf(logfile, "%s %s %s ERROR %s",
		time.Now().UTC().Format("2006-01-02T15:04:05.999999Z"), extensionShortName, version, message)

	fmt.Fprintf(logWriter, "%s\n", logMesssage)
}

func main() {
	// Manage closing of the file automatically once main returns
	initLogging()
	defer logWriter.Close()

	// Command line flags that are currently supported
	commandStringPtr := flag.String("command", "", "Desired command (install, enable, update, disable, uninstall")
	parseJSONPtr := flag.String("jsonfile", "", "Path to the JSON file loction")

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
	default:
		errorLogger.Printf("Command \"%s\" not recognized", *commandStringPtr)
	}

	// Parse the provided JSON file if there is one
	if *parseJSONPtr != "" {
		parseJSON(*parseJSONPtr)
	}
}
