package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/go-kit/kit/log"
)

var (
	version = "1.0.0.0"
)

func install() {
	operation := "install"
	msg := "Installed Successfully."

	// Configuring logger to print time and verison by default
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "time", log.DefaultTimestamp, "version", version)

	logger.Log("event", operation, "message", msg)

	/*	Format: time=XX version=XX event=XX message=XX
		Example: time=2020-XX-27T10:14:34.9357863-07:00 version=1.0.0.0 event=install message="Installed Successfully."
	*/
}

func enable() {
	fmt.Println("Enabled Successfully.")
}

func disable() {
	fmt.Println("Disabled Successfully.")
}

func uninstall() {
	fmt.Println("Uninstalled Successfully.")
}

func update() {
	fmt.Println("Updated Successfully.")
}

// Can be called with jsonfile=test (for example if you have a test.json file)
// TODO: Make accepting "jsonfile=test.json" possible
func parseJSON(jsonFilePath string) {
	//	Open the provided file
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println("File Not Found")
		os.Exit(1)
	}
	fmt.Println("File opened successfully")

	// Defer file closing until parseJSON() returns to its caller
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

		fmt.Println(key, reverseValue)
	}

}

func main() {
	if len(os.Args[1:]) > 0 {
		for _, a := range os.Args[1:] {
			/*	TODO : Not sure if there is a better method in regexp so don't need multiple vars
			 */
			matchDisable, _ := regexp.MatchString("^([-/]*)(disable)", a)
			matchUninstall, _ := regexp.MatchString("^([-/]*)(uninstall)", a)
			matchInstall, _ := regexp.MatchString("^([-/]*)(install)", a)
			matchEnable, _ := regexp.MatchString("^([-/]*)(enable)", a)
			matchUpdate, _ := regexp.MatchString("^([-/]*)(update)", a)

			if matchDisable {
				disable()
			} else if matchUninstall {
				uninstall()
			} else if matchInstall {
				install()
			} else if matchEnable {
				enable()
			} else if matchUpdate {
				update()
			} else {
				matchJSON, _ := regexp.MatchString("^([-/]*)(jsonfile=)", a)

				// This is a workaround for when "." is included in the command line args, it separates the args
				// TODO: implement this in a smarter way
				if matchJSON {
					s := a[10:] + ".json"
					parseJSON(s)
				}
			}
		}
	} else {
		fmt.Println("No command line arguments provided")
	}
	/* 	TODO : Error handling might be necessary for if there is no match, but this could
	just be a print statement else case if the regexp doesn't raise panics/errors
	*/
}
