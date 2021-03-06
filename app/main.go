package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"strconv"

	"github.com/Useurmind/spas/handler"
)

func main() {
	if sliceContains(os.Args, "-h") || sliceContains(os.Args, "--help") {
		printHelp()
		return
	}

	options, err := handler.GetOptions()
	if err != nil {
		panic(err)
	}

	SPASHandler := handler.NewSPASHandler(options)

	listenOn := fmt.Sprintf("%s:%s", options.Address, options.Port)
	log.Println("SPA server listening on", listenOn)
	if !options.ForceHTTP {
		log.Println("Serving https...")
		err = http.ListenAndServeTLS(listenOn, options.CertFilePath, options.KeyFilePath, SPASHandler)
	} else {
		log.Println("ForceHTTP set, serving http...")
		err = http.ListenAndServe(listenOn, SPASHandler)
	}
	if err != nil {
		panic(err)
	}
}

func printHelp() {
	_, appName := path.Split(strings.ReplaceAll(os.Args[0], "\\", "/"))
	defaultOptions, err := handler.DefaultOptions()
	if err != nil {
		panic(err)
	}

	fmt.Println("Usage:")
	fmt.Println(" ")
	fmt.Println("	", appName, "[options]")
	fmt.Println(" ")
	fmt.Println("Available options:")
	fmt.Println(" ")
	printOption("configfile", "	a path to a config file that contains the configuration for the spa server", defaultOptions.ConfigFile)
	printOption("address", "	address to listen on", defaultOptions.Address)
	printOption("port", "		port to listen on", defaultOptions.Port)
	printOption("servefolder", "	the folder to serve", "current working directory, e.g.: "+defaultOptions.ServeFolder)
	printOption("htmlindexfile", "path to the root index file of the spa app", defaultOptions.HTMLIndexFile)
	printOption("certfilepath", "	path to the ssl certificate (chain) for the server", defaultOptions.CertFilePath)
	printOption("keyfilepath", "	path to the private key of the ssl certificate of the server", defaultOptions.KeyFilePath)
	printOption("forcehttp", "	Set this flag to force serving http", strconv.FormatBool(defaultOptions.ForceHTTP))
}

func printOption(name string, description string, defaultValue string) {
	fmt.Println("	--"+name, description, "(default: "+defaultValue+")")
}

func sliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if value == v {
			return true
		}
	}

	return false
}
