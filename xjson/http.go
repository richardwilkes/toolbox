package xjson

import (
	"bytes"
	"net/http"
)

// GetRequest calls http.Get with the URL and returns the response body as a
// new JSON object.
func GetRequest(url string) (statusCode int, body *JSON, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err == nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		statusCode = resp.StatusCode
		body = MustParseJSONStream(resp.Body)
	} else {
		body = &JSON{}
	}
	return
}

// PostRequest calls http.Post with the URL and the contents of this JSON
// object and returns the response body as a new JSON object.
func (j *JSON) PostRequest(url string) (statusCode int, body *JSON, err error) {
	var resp *http.Response
	if resp, err = http.Post(url, "application/json", bytes.NewBuffer(j.Bytes())); err == nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		statusCode = resp.StatusCode
		body = MustParseJSONStream(resp.Body)
	} else {
		body = &JSON{}
	}
	return
}
