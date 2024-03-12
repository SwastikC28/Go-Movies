package middleware

import (
	"log"
	"net/http"
)

func ReqLogger(next http.Handler) http.Handler {
	// We wrap our anonymous function, and cast it to a http.HandlerFunc
	// Because our function signature matches ServeHTTP(w, r), this allows
	// our function (type) to implicitly satisify the http.Handler interface.
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Logic before - reading request values, putting things into the
			// request context, performing authentication

			// Print Req URL
			log.Println("REQ - ", r.Method, r.URL)

			// Important that we call the 'next' handler in the chain. If we don't,
			// then request handling will stop here.
			next.ServeHTTP(w, r)
			// Logic after - useful for logging, metrics, etc.
			//
			// It's important that we don't use the ResponseWriter after we've called the
			// next handler: we may cause conflicts when trying to write the response
		},
	)
}
