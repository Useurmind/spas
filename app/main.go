package main

import (
	"strings"
	"path"
	"os"
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, req *http.Request) {

}

func main() {
	if sliceContains(os.Args, "-h") || sliceContains(os.Args, "--help"){
		printHelp()
		return
	}

	options, err := GetOptions()
	if err != nil {
		panic(err)
	}

	SPASHandler := NewSPASHandler(options)

	listenOn := fmt.Sprintf("%s:%s", options.Address, options.Port)
	fmt.Println("spas listening on", listenOn)
	err = http.ListenAndServe(listenOn, SPASHandler)
	if err != nil {
		panic(err)
	}
}

func printHelp() {
	_, appName := path.Split(strings.ReplaceAll(os.Args[0], "\\", "/"))
	defaultOptions, err := DefaultOptions()
	if err != nil {
		panic(err)
	}

	fmt.Println("Usage:")
	fmt.Println(" ")
	fmt.Println("	", appName ,"[options]")
	fmt.Println(" ")
	fmt.Println("Available options:")
	fmt.Println(" ")
	printOption("configfile", "	a path to a config file that contains the configuration for the spa server", defaultOptions.ConfigFile)
	printOption("address", "	address to listen on", defaultOptions.Address)
	printOption("port", "		port to listen on", defaultOptions.Port)
	printOption("servefolder", "	the folder to serve", "current working directory, e.g.: " + defaultOptions.ServeFolder)
	printOption("htmlindexfile", "path to the root index file of the spa app", defaultOptions.HTMLIndexFile)
}

func printOption(name string, description string, defaultValue string) {
	fmt.Println("	--"+name, description, "(default: " + defaultValue + ")")
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
        if value == v {
            return true
        }   
	}
	
	return false
}