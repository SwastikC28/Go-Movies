package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
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

// UnmarshalForm parses form data into a struct.
func UnmarshalForm(values url.Values, v interface{}) error {
	var ErrInvalidType = errors.New("invalid type")

	// Get the type of the struct
	valueType := reflect.TypeOf(v)
	if valueType.Kind() != reflect.Ptr || valueType.Elem().Kind() != reflect.Struct {
		return ErrInvalidType
	}

	// Get the value of the struct
	value := reflect.ValueOf(v).Elem()

	// Iterate through the struct fields
	for i := 0; i < value.NumField(); i++ {
		// Get the field type and name
		fieldType := value.Type().Field(i)
		fieldName := fieldType.Name

		// Get the corresponding form value
		fieldValue := values.Get(fieldName)

		// Set the field value if the form value exists
		if fieldValue != "" {
			// Get the field value and set it
			field := value.Field(i)
			if field.CanSet() {
				// Convert the form value to the field type
				fieldKind := field.Kind()
				switch fieldKind {
				case reflect.String:
					field.SetString(fieldValue)
					// Add cases for other types as needed
				}
			}
		}
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

// Function to extract boundary parameter from Content-Type header
func ExtractBoundary(contentType string) string {
	parts := strings.Split(contentType, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "boundary=") {
			return strings.TrimPrefix(part, "boundary=")
		}
	}
	return ""
}
