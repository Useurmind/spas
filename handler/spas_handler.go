package handler

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
	err := h.ensureInit()
	if err != nil {
		w.WriteHeader(500)
		return
	}

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

func (h *SPASHandler) ensureInit() error {
	h.ensureFileHandler()
	err := h.ensureStaticFilesMap()
	if err != nil {
		return err
	}

	return nil
}

func (h *SPASHandler) ensureFileHandler() {
	if h.fileServerHandler == nil {
		h.fileServerHandler = http.FileServer(http.Dir(h.options.ServeFolder))
	}
}

func (h *SPASHandler) ensureStaticFilesMap() error {
	if len(h.staticFiles) == 0 {
		serveFolder, err := cleanFilePath(h.options.ServeFolder)
		if err != nil {
			return err
		}
	
		log.Println("Absolute serve folder is", serveFolder)

		err = filepath.Walk(serveFolder,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				cleanedPath, err := cleanFilePath(path)
				if err != nil {
					return err
				}

				if !info.IsDir() {
					// remember static file
					relativePath := strings.TrimPrefix(cleanedPath, serveFolder)
					urlPath := strings.ReplaceAll(relativePath, "\\", "/")
					urlPath = cleanURLPath(urlPath)
					h.staticFiles[urlPath] = cleanedPath

					log.Printf("Found static file %s, available under url path %s", cleanedPath, urlPath)
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	}

	return nil
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

// make sure input and file paths match
func cleanFilePath(filePath string) (string, error) {
	cleanedPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	return cleanedPath, nil
}

// remove double quotes
// insert initial quote if missing
func cleanURLPath(urlPath string) string {
	cleanedPath := fmt.Sprintf("/%s", urlPath)
	cleanedPath = strings.ReplaceAll(cleanedPath, "//", "/")

	return cleanedPath
}
