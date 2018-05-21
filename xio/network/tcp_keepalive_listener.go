package network

import (
	"net"
	"time"
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted connections
// so that dead TCP connections (e.g. closing laptop mid-download) eventually
// go away. This code was mostly copied from the http package, with changes to
// make it publicly accessible and to check some errors that were ignored.
type TCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements the Accept method in the Listener interface; it
// waits for the next call and returns a generic Conn.
func (listener TCPKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err = tc.SetKeepAlive(true); err != nil {
		return nil, err
	}
	if err = tc.SetKeepAlivePeriod(3 * time.Minute); err != nil {
		return nil, err
	}
	return tc, nil
}
