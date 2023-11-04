package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

// ... other code

func isValidExtension(filePath string, allowedExtensions []string) bool {
	extension := filepath.Ext(filePath)
	for _, allowedExtension := range allowedExtensions {
		if strings.EqualFold(extension, allowedExtension) {
			return true
		}
	}
	return false
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	allowedExtensions := []string{".html", ".css", ".js", ".ico"}
	if !isValidExtension(r.URL.Path, allowedExtensions) {
		http.NotFound(w, r)
		return
	}

	fs := http.FileServer(http.Dir("."))
	fs.ServeHTTP(w, r)
}
