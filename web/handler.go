package web

import (
	"net/http"
	// "os"
	// "path/filepath"
	"text/template"
)

type renderData struct {
	ApplicationName string
}

type SpaHandler struct {
	StaticPath string
	IndexPath  string
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	view, _ := template.ParseFiles("build/static/index.html")
	view.Execute(w, renderData{
		ApplicationName: "GoExpense Built",
	})
}

// func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// get the absolute path to prevent directory traversal
// 	path, err := filepath.Abs(r.URL.Path)
// 	if err != nil {
// 		// if we failed to get the absolute path respond with a 400 bad request
// 		// and stop
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// prepend the path with the path to the static directory
// 	path = filepath.Join(h.StaticPath, path)

// 	// check whether a file exists at the given path
// 	_, err = os.Stat(path)
// 	if os.IsNotExist(err) {
// 		// file does not exist, serve index.html
// 		view, _ := template.ParseFiles("dist/index.html")
// 		view.Execute(w, renderData{
// 			ApplicationName: "GoExpense Built",
// 		})
// 		// http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
// 		return
// 	} else if err != nil {
// 		// if we got an error (that wasn't that the file doesn't exist) stating the
// 		// file, return a 500 internal server error and stop
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// otherwise, use http.FileServer to serve the static dir
// 	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
// }
