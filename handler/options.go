package handler

import (
	"encoding/json"
	"path/filepath"
	"strconv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
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

	// The path to the file with the server ssl certificate (chain) to use.
	CertFilePath string

	// The path to the file with the private key for the server ssl certificate.
	KeyFilePath string

	// If this is set the spas server will serve http.
	ForceHTTP bool
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
	options.CertFilePath = "spas_cert.pem"
	options.KeyFilePath = "spas_key.pem"
	options.ForceHTTP = false

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

	certFilePath, err := getOption("CertFilePath", options.CertFilePath, defaultOptions.CertFilePath)
	if err != nil {
		return nil, err
	}

	keyFilePath, err := getOption("KeyFilePath", options.KeyFilePath, defaultOptions.KeyFilePath)
	if err != nil {
		return nil, err
	}

	forceHTTP, err := getOptionFlag("ForceHTTP", options.ForceHTTP, defaultOptions.ForceHTTP)
	if err != nil {
		return nil, err
	}

	options.ConfigFile = configFile
	options.Address = address
	options.Port = port
	options.ServeFolder = serveFolder
	options.HTMLIndexFile = htmlIndexFile
	options.CertFilePath = certFilePath
	options.KeyFilePath = keyFilePath
	options.ForceHTTP = forceHTTP

	log.Println("Applied options are:")
	options.Log()

	err = options.WarnProblems()
	if err != nil {
		return nil, err
	}

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

// WarnProblems writes problems in the options to the log output.
func (o *Options) WarnProblems() error {
	if o.CertFilePath != "" && o.KeyFilePath == "" {
		log.Println("WARNING: CertFilePath is set, but KeyFilePath is missing. If you want to serve SSL you need a key file containing the private key for the SSL certificate")
	}

	absServeFolder, err := filepath.Abs(o.ServeFolder)
	if err != nil {
		return err
	}

	if o.KeyFilePath != "" {
		keyFileDirectory, _ := path.Split(o.KeyFilePath)

		absKeyFileDir, err := filepath.Abs(keyFileDirectory)
		if err != nil {
			return err
		}

		if strings.HasPrefix(absKeyFileDir, absServeFolder) {
			log.Println("WARNING: Private ssl key lies in", absKeyFileDir, "which is a subdirectory of the ServeFolder", absServeFolder, ". This is very dangerous as the private ssl key will be served.")
		}
	}

	if absServeFolder == "/" {
		log.Println("WARNING: You are serving the root folder. This is very dangerous as the whole file system will be served. Put your served files in a separate folder.")
	}

	return nil
}

func getOptionFlag(optionName string, configFileValue bool, defaultValue bool) (bool, error) {
	strValue, err := getOption(optionName, strconv.FormatBool(configFileValue), strconv.FormatBool(defaultValue))

	if err != nil {
		// the only error ever return is when no value is found for an option on the command line
		// this indicates that the flag is set
		return true, nil
	}

	return strconv.ParseBool(strValue)
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
		value, found = os.LookupEnv(fmt.Sprintf("SPAS_%s", varName))
	}

	if !found && configFileValue != "" {
		log.Printf("  Applying option from config file")
		value = configFileValue
	}

	if !found && value == "" {
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

		if arg == optionName {
			if i < len(os.Args)-1 {
				argValue := os.Args[i+1]
				return true, argValue, nil
			}

			return true, "", fmt.Errorf("No value given for option %s", optionName)
		}
	}

	return false, "", nil
}
