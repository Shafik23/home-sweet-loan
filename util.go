package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

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
	// fmt.Println("Requested path", r.URL.Path)

	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./hsl.html")
		return
	}

	allowedExtensions := []string{".html", ".css", ".js", ".ico", ".png", ".jpg", ".jpeg", ".svg", ".woff2", ".woff", ".ttf", ".eot"}
	if !isValidExtension(r.URL.Path, allowedExtensions) {
		http.NotFound(w, r)
		return
	}

	fs := http.FileServer(http.Dir("."))
	fs.ServeHTTP(w, r)
}
