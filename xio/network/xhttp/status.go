package xhttp

import (
	"net/http"
)

// WriteHTTPStatus sends an HTTP response header with 'statusCode' and follows
// it with the standard text for that code as the body.
func WriteHTTPStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	// The extra code here is just to quiet the linter about not checking
	// for an error.
	if _, err := w.Write([]byte(http.StatusText(statusCode))); err != nil {
		return
	}
}
