package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"shared/security"
	"strings"
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

func AccessGuard(next http.HandlerFunc, isAdmin bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			RespondJSON(w, http.StatusUnauthorized, "User unauthorized to access this route")
			return
		}

		// Extract Token from Authorization Header
		token = strings.Split(token, " ")[1]

		fmt.Println("TOKEN:", token)

		// Verify Token
		payload, err := security.Verify(token)
		if err != nil {
			fmt.Println(err)
			RespondJSON(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// Check if admin is required
		if isAdmin && !payload.IsAdmin {
			RespondJSON(w, http.StatusUnauthorized, "User unauthorized to access this route")
			return
		}

		ctx := payload.ToContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
