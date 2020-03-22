package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"log"
)

// Options for the spa server.
// The options will be applied in the following precedence (top most wins)
// 1. command line arguments
// 2. environment variables
// 3. ConfigFile 
// 
// Command line arguments are added with all lowercase names of the actual options preceeded by a double slash.
// __Example__: `spas --port 8090`
//
// Environment variables are named all uppercase with a prefix of SPAS_
// __Example__: export SPAS_PORT=8090
//
// The config file is either `spas.config.json` in the current working directory or you need to specify it either via 
// command line or environment variable.
type Options struct {
	// A config file that is used for configuration.
	// The config file is optional.
	// Default spas.config.json in the working directory.
	ConfigFile string

	// The address to listen on.
	// Example:
	// - 127.0.0.1
	// Empty by default, will listen on all.
	Address string

	// The port to listen on. The default is 8080.
	Port string

	// The folder which should be served, default current working directory.
	// All files found in this folder will be served.
	// Even files in subfolders.
	ServeFolder string

	// The path/name of the HTML index file that is the root of your app.
	// The index file usually loads the root js bundle file for rendering the SPA.
	// Default is index.html.
	HTMLIndexFile string
}

// DefaultOptions returns the default options.
func DefaultOptions() (*Options, error) {
	options := Options{}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	options.ConfigFile = "spas.config.json"
	options.Address = ""
	options.Port = "8080"
	options.ServeFolder = wd
	options.HTMLIndexFile = "index.html"

	return &options, nil
}

// GetOptions gets the options for the server.
// It will read command line and environment variables.
// On command line the pattern is: --<lowerCaseSettingName> <value>
// For environment variables the pattern is: SPAS_<UPPERCASESETTINGNAME>
// Command line values overwrite environment variables.
func GetOptions() (*Options, error) {
	defaultOptions, err := DefaultOptions()
	if err != nil {
		return nil, fmt.Errorf("Could not get default options, %s", err)
	}
	options := Options{}

	configFile, err := getOption("configfile", "", "spas.config.json")
	if err != nil {
		return nil, err
	}

	if configFile != "" {
		configFileContent, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("Could not find config file in %s (%s), skipping", configFile, err)
		} else {

			err := json.Unmarshal(configFileContent, &options)
			if err != nil {
				log.Panicf("Could not read config file %s: %s", configFile, err)
			}

		}
	}	

	address, err := getOption("address", options.Address, defaultOptions.Address)
	if err != nil {
		return nil, err
	}

	port, err := getOption("port", options.Port, defaultOptions.Port)
	if err != nil {
		return nil, err
	}

	serveFolder, err := getOption("servefolder", options.ServeFolder, defaultOptions.ServeFolder)
	if err != nil {
		return nil, err
	}

	htmlIndexFile, err := getOption("HTMLIndexFile", options.HTMLIndexFile, defaultOptions.HTMLIndexFile)
	if err != nil {
		return nil, err
	}

	options.ConfigFile = configFile
	options.Address = address
	options.Port = port
	options.ServeFolder = serveFolder
	options.HTMLIndexFile = htmlIndexFile

	log.Println("Applied options are:")
	options.Log()

	return &options, nil
}

// Log the options as a json object to the log output.
func (o *Options) Log() {
	jsonOptions, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Printf("ERROR: Could not print options, %s", err)
	} else {
		log.Printf(string(jsonOptions))
	}
}

func getOption(optionName string, configFileValue string, defaultValue string) (string, error) {
	log.Printf("Searching option %s (defaultValue is %s)", optionName, defaultValue)
	found, value, err := findArgForOption(fmt.Sprintf("--%s", strings.ToLower(optionName)))
	if err != nil {
		return "", err
	}

	if !found {
		varName := strings.ToUpper(optionName)
		log.Printf("  Trying to get option from environment variable %s", varName)
		value = os.Getenv(fmt.Sprintf("SPAS_%s", varName))
	}

	if value == "" && configFileValue != "" {
		log.Printf("  Applying option from config file")
		value = configFileValue
	}

	if value == "" {
		log.Printf("  Option nowhere found, applying default value %s", defaultValue)
		value = defaultValue
	} else {
		log.Printf("  Option set to value %s", value)
	}

	return value, nil
}

func findArgForOption(optionName string) (bool, string, error) {	
	log.Printf("  Trying to get option from command line argument %s", optionName)
	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]

		if arg == optionName  {
			if i < len(os.Args) - 1 {
				argValue := os.Args[i+1]
				return true, argValue, nil
			} 
			
			return true, "", fmt.Errorf("No value given for option %s", optionName)
		}
	}

	return false, "", nil
}