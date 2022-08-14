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
			// if we failed to get the absolute path respond with a 400 bad request
			// and stop
			http.Error(w, err.Error(), http.StatusBadRequest)
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
			// if we got an error (that wasn't that the file doesn't exist) stating the
			// file, return a 500 internal server error and stop
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// otherwise, use http.FileServer to serve the static dir
		http.FileServer(http.Dir(staticPath)).ServeHTTP(w, r)
	}
}
