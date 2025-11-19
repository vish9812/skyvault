package api

import (
	"net/http"
	"os"
	"path/filepath"
)

// spaHandler serves static files and falls back to index.html for SPA routes.
// This enables client-side routing for Single Page Applications.
//
// How it works:
// - If the requested file exists (e.g., /assets/main.js), serve it normally
// - If the file doesn't exist (e.g., /drive/folder-id), serve index.html
// - This allows the SPA router to handle the route on the client side
func spaHandler(staticDir string) http.HandlerFunc {
	fileServer := http.FileServer(http.Dir(staticDir))

	return func(w http.ResponseWriter, r *http.Request) {
		// Build the full file path
		path := filepath.Join(staticDir, r.URL.Path)

		// Check if file exists
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// File does not exist, serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		} else if err != nil {
			// Some other error occurred
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// File exists, serve it normally
		fileServer.ServeHTTP(w, r)
	}
}
