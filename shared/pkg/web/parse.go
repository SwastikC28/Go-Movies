package web

import (
	"net/http"
	"strings"
)

// Get all associations from the request
func ParseAssociation(r *http.Request) []string {
	params := r.URL.Query()

	associations := params.Get("includes")

	if associations == "" {
		return []string{}
	}

	return strings.Split(associations, ",")
}
