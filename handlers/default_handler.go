package handlers

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body := make(map[string]string)
	body["message"] = "Ok"
	jData, _ := json.Marshal(body)
	w.Write(jData)
}
