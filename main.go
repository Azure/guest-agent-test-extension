package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

var (
	version            = "1.0.0.0"
	extensionShortName = "GATestExt"
	// Logging is currently set up to create/add to the logile in the directory from where the binary is executed
	logfile = "guest-agent-test-extension.log"

	infoLogger  *log.Logger
	errorLogger *log.Logger
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
func initLogging() *os.File {
	// Open the file as read only
	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		errorLogger.Println("Error opening file: ", logfile)
		log.Fatal(err)
	}

	// Format the timestamp to match the UTC format from the waagent log files
	// Sample out: 2020-07-29T01:10:53.960425Z GATestExt version: 1.0.0.0 INFO main.go:22: Hello World!
	infoLogger = log.New(file, "INFO ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	// Sample out: 2020-07-29T01:10:53.960425Z GATestExt version: 1.0.0.0 ERROR main.go:22: Hello World!
	errorLogger = log.New(file, "ERROR ", log.LUTC|log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	return file
}

func main() {
	// Manage closing of the file automatically once main returns
	file := initLogging()
	defer file.Close()

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
