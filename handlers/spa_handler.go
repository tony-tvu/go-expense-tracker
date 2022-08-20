package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

func SpaHandler(staticPath string, indexPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the absolute path to prevent directory traversal
		path, err := filepath.Abs(r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// prepend the path with the path to the static directory
		path = filepath.Join(staticPath, path)

		// check whether a file exists at the given path
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// otherwise, use http.FileServer to serve the static dir
		http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
	}
}
