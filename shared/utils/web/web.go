package web

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, statusCode int, value interface{}) {
	// Marshal into JSON
	body, err := json.MarshalIndent(value, "", "\t")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Set Header
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}
