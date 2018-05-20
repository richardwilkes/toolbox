package xhttp

import (
	"fmt"
	"net/http"
)

// WriteHTTPStatus sends an HTTP response header with 'statusCode' and follows
// it with the standard text for that code as the body.
func WriteHTTPStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, http.StatusText(statusCode))
}
