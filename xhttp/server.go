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
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/richardwilkes/toolbox/v2/xos"
	"github.com/richardwilkes/toolbox/v2/xreflect"
)

// DefaultStopGracePeriod is the default maximum time to wait for the server to stop when Shutdown is called.
const DefaultStopGracePeriod = 5 * time.Second

// ServerConfig holds configuration for an HTTP server.
type ServerConfig struct {
	// Handler is the root HTTP handler to use.
	Handler http.Handler
	// Logger is the logger to use. If nil, slog.Default() will be used.
	Logger *slog.Logger
	// TLSConfig optionally provides a TLS configuration for use by ServeTLS and ListenAndServeTLS. Note that this value
	// is cloned by ServeTLS and ListenAndServeTLS, so it's not possible to modify the configuration with methods like
	// tls.Config.SetSessionTicketKeys.
	TLSConfig *tls.Config
	// CertFile optionally provides a cert file to use for TLS connections. Must be paired with a KeyFile.
	CertFile string
	// KeyFile optionally provides a key file to use for TLS connections. Must be paired with a CertFile.
	KeyFile string
	// Host is the host name or IP address of the server. If empty, all addresses of the host will be used.
	Host string
	// Port is the port to listen on. A value outside the range 1-65535 will result in a dynamically assigned port.
	Port int
	// ReadTimeout is the maximum duration for reading the entire request, including the body. A zero or negative value
	// means there will be no timeout.
	//
	// Because ReadTimeout does not let Handlers make per-request decisions on each request body's acceptable deadline
	// or upload rate, most users will prefer to use ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration
	// ReadHeaderTimeout is the amount of time allowed to read request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what is considered too slow for the body. If zero, the value
	// of ReadTimeout is used. If negative, or if zero and ReadTimeout is zero or negative, there is no timeout.
	ReadHeaderTimeout time.Duration
	// WriteTimeout is the maximum duration before timing out writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not let Handlers make decisions on a per-request basis. A
	// zero or negative value means there will be no timeout.
	WriteTimeout time.Duration
	// IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled. If zero, the
	// value of ReadTimeout is used. If negative, or if zero and ReadTimeout is zero or negative, there is no timeout.
	IdleTimeout time.Duration
	// StopGracePeriod is the maximum time to wait for the server to stop. If zero or negative, the default of
	// DefaultStopGracePeriod is used.
	StopGracePeriod time.Duration
	// MaxHeaderBytes controls the maximum number of bytes the server will read parsing the request header's keys and
	// values, including the request line. It does not limit the size of the request body. If zero,
	// http.DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int
	// DisableGeneralOptionsHandler, if true, passes "OPTIONS *" requests to the Handler, otherwise responds with 200 OK
	// and Content-Length: 0.
	DisableGeneralOptionsHandler bool
}

// Server manages a http server.
type Server struct {
	server          *http.Server
	originalHandler http.Handler
	logger          *slog.Logger
	listenerAddress net.Addr
	err             error
	started         chan struct{}
	stopped         chan struct{}
	protocol        string
	certFile        string
	keyFile         string
	lock            sync.Mutex
	stopID          int
	stopGracePeriod time.Duration
}

// NewServer creates a new http service. The server is not started until you call Run.
func NewServer(cfg *ServerConfig) (*Server, error) {
	port := cfg.Port
	if port < 1 || port > 65535 {
		port = 0
	}
	s := Server{
		server: &http.Server{
			Addr:                         net.JoinHostPort(cfg.Host, strconv.Itoa(port)),
			DisableGeneralOptionsHandler: cfg.DisableGeneralOptionsHandler,
			TLSConfig:                    cfg.TLSConfig,
			ReadTimeout:                  cfg.ReadTimeout,
			ReadHeaderTimeout:            cfg.ReadHeaderTimeout,
			WriteTimeout:                 cfg.WriteTimeout,
			IdleTimeout:                  cfg.IdleTimeout,
			MaxHeaderBytes:               cfg.MaxHeaderBytes,
		},
		originalHandler: cfg.Handler,
		logger:          cfg.Logger,
		certFile:        cfg.CertFile,
		keyFile:         cfg.KeyFile,
		started:         make(chan struct{}),
		stopped:         make(chan struct{}),
		stopGracePeriod: cfg.StopGracePeriod,
	}
	if s.stopGracePeriod <= 0 {
		s.stopGracePeriod = DefaultStopGracePeriod
	}
	if (cfg.CertFile != "" && cfg.KeyFile != "") ||
		(cfg.TLSConfig != nil &&
			(len(cfg.TLSConfig.Certificates) > 0 ||
				cfg.TLSConfig.GetCertificate != nil ||
				cfg.TLSConfig.GetConfigForClient != nil)) {
		s.protocol = "https"
	} else {
		s.protocol = "http"
	}
	if s.logger == nil {
		s.logger = slog.Default()
	}
	if xreflect.IsNil(s.originalHandler) {
		s.originalHandler = http.NewServeMux()
	}
	s.server.Handler = &s
	return &s, nil
}

// Protocol returns the protocol this server is handling.
func (s *Server) Protocol() string {
	return s.protocol
}

// Run the server. It blocks until the server stops or an error occurs.
func (s *Server) Run() error {
	s.lock.Lock()
	started := s.stopID != 0
	if started {
		s.lock.Unlock()
		return errors.New("server has already been started")
	}
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		s.lock.Unlock()
		return err
	}
	s.listenerAddress = listener.Addr()
	s.stopID = xos.RunAtExit(s.Stop)
	defer func() {
		s.lock.Lock()
		xos.CancelRunAtExit(s.stopID)
		s.stopID = 0
		s.lock.Unlock()
		close(s.stopped)
	}()
	slog.Info(s.protocol + " server is now listening on " + s.ListenerAddress().String())
	close(s.started)
	s.lock.Unlock()
	if s.protocol == "https" {
		err = s.server.ServeTLS(listener, s.certFile, s.keyFile)
	} else {
		err = s.server.Serve(listener)
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.lock.Lock()
		s.err = err
		s.lock.Unlock()
	}
	return nil
}

// ListenerAddress returns the address of the underlying listener. This will return nil until the server is started.
func (s *Server) ListenerAddress() net.Addr {
	return s.listenerAddress
}

// Started returns true if the server is started.
func (s *Server) Started() bool {
	select {
	case <-s.started:
		return true
	default:
		return false
	}
}

// WaitForStart blocks until the server starts. If the server is already started, it returns immediately.
func (s *Server) WaitForStart() {
	<-s.started
}

// Stopped returns true if the server is stopped.
func (s *Server) Stopped() bool {
	select {
	case <-s.stopped:
		return true
	default:
		return false
	}
}

// WaitForStop blocks until the server stops. If the server is already stopped, it returns immediately.
func (s *Server) WaitForStop() {
	<-s.stopped
}

// Error returns any error that occurred while running the server.
func (s *Server) Error() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.err
}

// Stop the server. It blocks until the server stops or the grace period is reached. If you call this from within a
// handler, you should call it asynchronously since the server attempts to allow existing requests to complete.
func (s *Server) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.stopID == 0 {
		return
	}
	address := s.ListenerAddress().String()
	slog.Info("stopping " + s.protocol + " server listening on " + address)
	defer func() {
		xos.CancelRunAtExit(s.stopID)
		s.stopID = 0
		slog.Info(s.protocol + " server listening on " + address + " has stopped")
	}()
	ctx, cancel := context.WithTimeout(context.Background(), s.stopGracePeriod)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.err = err
	}
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := s.logger.With(
		"method", req.Method,
		"url", req.URL.String(),
		"client-ip", ClientIP(req),
	)
	md := &Metadata{Logger: logger}
	sw := NewStatusWriter(w, req)
	started := time.Now()
	defer func() {
		since := time.Since(started)
		millis := int64(since / time.Millisecond)
		micros := int64(since/time.Microsecond) - millis*1000
		var msg string
		if md.LogMsg == "" {
			msg = "request complete"
		} else {
			msg = md.LogMsg
		}
		logger.Info(msg,
			"status", sw.Status(),
			"bytes", sw.BytesWritten(),
			slog.String("elapsed", fmt.Sprintf("%d.%03dms", millis, micros)),
		)
	}()
	defer xos.PanicRecovery(func(err error) {
		logger.Error("recovered from panic in handler", "error", err)
		// It might be too late to set the header, but we'll try anyway.
		w.WriteHeader(http.StatusInternalServerError)
	})
	ctx := metadataInContext(req.Context(), md)
	if s.server.WriteTimeout > 0 {
		// Update the context to expire at the same time as the write timeout.
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.server.WriteTimeout)
		defer cancel()
	}
	s.originalHandler.ServeHTTP(sw, req.WithContext(ctx))
}
