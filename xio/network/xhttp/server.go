// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package xhttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/toolbox/v2/xio/network"
	"github.com/richardwilkes/toolbox/v2/xos"
)

// Constants for protocols the server can provide.
const (
	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"
)

type ctxKey int

const metadataKey ctxKey = 1

// Metadata holds auxiliary information for a request.
type Metadata struct {
	Logger *slog.Logger
	User   string
}

// Server holds the data necessary for the server.
type Server struct {
	WebServer           *http.Server
	Logger              *slog.Logger
	clientHandler       http.Handler
	StartedChan         chan struct{} // If not nil, will be closed once the server is ready to accept connections
	ShutdownCallback    func(*slog.Logger)
	CertFile            string
	KeyFile             string
	addresses           []string
	ShutdownGracePeriod time.Duration
	port                int
	shutdownID          int
}

// Protocol returns the protocol this server is handling.
func (s *Server) Protocol() string {
	if s.CertFile != "" && s.KeyFile != "" {
		return ProtocolHTTPS
	}
	return ProtocolHTTP
}

// Addresses returns the host addresses being listened to.
func (s *Server) Addresses() []string {
	return s.addresses
}

// Port returns the port being listened to.
func (s *Server) Port() int {
	return s.port
}

// LocalBaseURL returns the local base URL that will reach the server.
func (s *Server) LocalBaseURL() string {
	return fmt.Sprintf("%s://%s:%d", s.Protocol(), network.IPv4LoopbackAddress, s.port)
}

func (s *Server) String() string {
	var buffer strings.Builder
	buffer.WriteString(s.Protocol())
	buffer.WriteString(" on ")
	for i, addr := range s.addresses {
		if i != 0 {
			buffer.WriteString(", ")
		}
		fmt.Fprintf(&buffer, "%s:%d", addr, s.port)
	}
	return buffer.String()
}

// Run the server. Does not return until the server is shutdown.
func (s *Server) Run() error {
	s.shutdownID = xos.RunAtExit(s.Shutdown)
	if s.Logger == nil {
		s.Logger = slog.Default()
	}
	s.clientHandler = s.WebServer.Handler
	s.WebServer.Handler = s
	var listener net.Listener
	_, _, err := net.SplitHostPort(s.WebServer.Addr)
	if err == nil {
		listener, err = net.Listen("tcp", s.WebServer.Addr)
	} else {
		listener, err = net.Listen("tcp", net.JoinHostPort(s.WebServer.Addr, "0"))
	}
	if err != nil {
		return errs.Wrap(err)
	}
	var host, portStr string
	if host, portStr, err = net.SplitHostPort(listener.Addr().String()); err != nil {
		return errs.Wrap(err)
	}
	if s.port, err = strconv.Atoi(portStr); err != nil {
		return errs.Wrap(err)
	}
	s.addresses = network.AddressesForHost(host)
	s.Logger.Info("listening", "protocol", s.Protocol(), "addresses", s.addresses, "port", s.port)
	if s.StartedChan != nil {
		go func() { close(s.StartedChan) }()
	}
	if s.Protocol() == ProtocolHTTPS {
		err = s.WebServer.ServeTLS(listener, s.CertFile, s.KeyFile)
	} else {
		err = s.WebServer.Serve(listener)
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errs.Wrap(err)
	}
	return nil
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	started := time.Now()
	sw := &StatusResponseWriter{
		Original: w,
		Head:     r.Method == http.MethodHead,
	}
	md := &Metadata{Logger: s.Logger.With("method", r.Method, "url", r.URL)}
	r = r.WithContext(context.WithValue(r.Context(), metadataKey, md))
	defer func() {
		if recovered := recover(); recovered != nil {
			err, ok := recovered.(error)
			if !ok {
				err = errs.Newf("%+v", recovered)
			}
			errs.LogTo(md.Logger, errs.NewWithCause("recovered from panic in handler", err))
			ErrorStatus(sw, http.StatusInternalServerError)
		}
		since := time.Since(started)
		millis := int64(since / time.Millisecond)
		micros := int64(since/time.Microsecond) - millis*1000
		written := sw.BytesWritten()
		md.Logger.Info("web", "status", sw.Status(), "bytes", written, "elapsed",
			fmt.Sprintf("%d.%03dms", millis, micros))
	}()
	s.clientHandler.ServeHTTP(sw, r)
}

// Shutdown the server gracefully.
func (s *Server) Shutdown() {
	xos.CancelRunAtExit(s.shutdownID)
	startedAt := time.Now()
	logger := s.Logger.With("protocol", s.Protocol(), "addresses", s.addresses, "port", s.port)
	logger.Info("starting shutdown")
	defer func() { logger.Info("finished shutdown", "elapsed", time.Since(startedAt)) }()
	gracePeriod := s.ShutdownGracePeriod
	if gracePeriod <= 0 {
		gracePeriod = time.Minute
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(gracePeriod))
	defer cancel()
	if err := s.WebServer.Shutdown(ctx); err != nil {
		errs.LogTo(logger, errs.NewWithCause("unable to shutdown gracefully", err))
	}
	if s.ShutdownCallback != nil {
		s.ShutdownCallback(logger)
	}
}

// MetadataFromRequest returns the Metadata from the request.
func MetadataFromRequest(req *http.Request) *Metadata {
	if md, ok := req.Context().Value(metadataKey).(*Metadata); ok {
		return md
	}
	return nil
}
