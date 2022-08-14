package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		body := make(map[string]string)
		body["message"] = "Ok"
		jData, _ := json.Marshal(body)
		w.Write(jData)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
