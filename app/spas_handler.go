package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// SPASHandler is the custom spas http handler.
type SPASHandler struct {
	// options to configure the behaviour of spas
	options *Options

	// A map of static files found in ServeFolder
	// It maps from the url path to the file path.
	staticFiles map[string]string

	// The file server handler that serves the files for the
	// spa.
	fileServerHandler http.Handler
}

func NewSPASHandler(options *Options) SPASHandler {
	return SPASHandler{
		options: options,
		staticFiles: make(map[string]string),
	}
}

func (h SPASHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.ensureInit()

	filePath, hasStaticFile := h.staticFiles[req.URL.Path]

	if pathEndsOnExtension(req.URL.Path) && !hasStaticFile {
		// a static file is requested, but there is none, return 404
		w.WriteHeader(404)
		fmt.Fprintf(w, "File not found")
		return
	} else if !hasStaticFile {
		filePath = path.Join(h.options.ServeFolder, h.options.HTMLIndexFile)
	}

	log.Printf("Serving url %s with file %s", req.RequestURI, filePath)

	http.ServeFile(w, req, filePath)
}

func (h *SPASHandler) ensureInit() {
	h.ensureFileHandler()
	h.ensureStaticFilesMap()
}

func (h *SPASHandler) ensureFileHandler() {
	if h.fileServerHandler == nil {
		h.fileServerHandler = http.FileServer(http.Dir(h.options.ServeFolder))
	}
}

func (h *SPASHandler) ensureStaticFilesMap() {
	if len(h.staticFiles) == 0 {
		err := filepath.Walk(h.options.ServeFolder,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					// remember static file
					relativePath := strings.TrimPrefix(path, h.options.ServeFolder)
					urlPath := strings.ReplaceAll(relativePath, "\\", "/")
					urlPath = cleanURLPath(urlPath)
					h.staticFiles[urlPath] = path

					log.Printf("Found static file %s, available under url path %s", path, urlPath)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}

	}
}

// the urlPath should start with a single slash
func (h *SPASHandler) hasStaticFile(urlPath string) bool {
	_, ok := h.staticFiles[urlPath]

	return ok
}

func pathEndsOnExtension(urlPath string) bool {
	cleanPath := cleanURLPath(urlPath)
	// pathParts := strings.Split(cleanPath, "/")
	// possibleFileName := pathParts[len(pathParts) - 1]
	return path.Ext(cleanPath) != ""
}

// remove double quotes
// insert initial quote if missing
func cleanURLPath(urlPath string) string {
	cleanedPath := fmt.Sprintf("/%s", urlPath)
	cleanedPath = strings.ReplaceAll(cleanedPath, "//", "/")

	return cleanedPath
}
