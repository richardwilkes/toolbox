package web

import (
	"net/http"
	"strings"
)

type routeCtxKey int

var routeKey routeCtxKey = 1

type route struct {
	path string
}

func (r *route) shift() string {
	i := strings.Index(r.path[1:], "/") + 1
	if i <= 0 {
		p := r.path[1:]
		r.path = "/"
		return p
	}
	p := r.path[1:i]
	r.path = r.path[i:]
	return p
}

// PathHeadThenShift returns the head segment of the request's adjusted path,
// then shifts it left by one segment. This does not adjust the path stored in
// req.URL.Path.
func PathHeadThenShift(req *http.Request) string {
	r, ok := req.Context().Value(routeKey).(*route)
	if !ok {
		return req.URL.Path
	}
	return r.shift()
}
