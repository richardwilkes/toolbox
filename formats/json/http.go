package json

import (
	"bytes"
	"net/http"
)

// GetRequest calls http.Get with the URL and returns the response body as a
// new Data object.
func GetRequest(url string) (statusCode int, body *Data, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err == nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		statusCode = resp.StatusCode
		body = MustParseStream(resp.Body)
	} else {
		body = &Data{}
	}
	return
}

// PostRequest calls http.Post with the URL and the contents of this Data
// object and returns the response body as a new Data object.
func (j *Data) PostRequest(url string) (statusCode int, body *Data, err error) {
	var resp *http.Response
	if resp, err = http.Post(url, "application/json", bytes.NewBuffer(j.Bytes())); err == nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		statusCode = resp.StatusCode
		body = MustParseStream(resp.Body)
	} else {
		body = &Data{}
	}
	return
}
