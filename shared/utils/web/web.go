package web

import (
	"encoding/json"
	"errors"
	"io"
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

// Unmarshals Request's Body into out variable
func UnmarshalJSON(r *http.Request, out interface{}) error {
	// Check if request is empty
	if r.Body == nil {
		return errors.New("body is empty")
	}

	// Read Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Unmarshal JSON
	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}

	return nil
}
