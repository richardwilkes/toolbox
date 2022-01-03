// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package network

import (
	"net"
	"time"
)

// TCPKeepAliveListener sets TCP keep-alive timeouts on accepted connections so that dead TCP connections (e.g. closing
// laptop mid-download) eventually go away. This code was mostly copied from the http package, with changes to make it
// publicly accessible and to check some errors that were ignored.
type TCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept implements the Accept method in the Listener interface; it waits for the next call and returns a generic Conn.
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
