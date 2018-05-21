package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/log/logadapter"
	"github.com/richardwilkes/toolbox/xio/network"
	"github.com/richardwilkes/toolbox/xio/network/xhttp"
)

// Constants for protocols the server can provide.
const (
	ProtocolHTTP  = "http"
	ProtocolHTTPS = "https"
)

// Server holds the data necessary for the server.
type Server struct {
	CertFile            string
	KeyFile             string
	ShutdownGracePeriod time.Duration
	Logger              logadapter.Logger
	WebServer           *http.Server
	addresses           []string
	port                int
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
	atexit.Register(s.Shutdown)
	if s.Logger == nil {
		s.Logger = &logadapter.Discarder{}
	}
	handler := s.WebServer.Handler
	s.WebServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		started := time.Now()
		sw := &xhttp.StatusResponseWriter{Original: w}
		defer func() {
			if err := recover(); err != nil {
				s.Logger.Error(errs.Newf("recovered from panic in handler\n%+v", err))
				sw.WriteHeader(http.StatusInternalServerError)
			}
			s.Logger.Infof("%d | %sms | %s bytes | %s %s", sw.Status(), humanize.Comma(int64(time.Since(started)/time.Millisecond)), humanize.Comma(int64(sw.BytesWritten())), req.Method, req.URL)
		}()
		handler.ServeHTTP(sw, req)
	})
	address := s.WebServer.Addr
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		address = address + ":0"
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return errs.Wrap(err)
	}
	listener := network.TCPKeepAliveListener{TCPListener: ln.(*net.TCPListener)}
	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		return errs.Wrap(err)
	}
	if s.port, err = strconv.Atoi(portStr); err != nil {
		return errs.Wrap(err)
	}
	s.addresses = network.AddressesForHost(host)
	s.Logger.Infof("Listening for %v", s)
	if s.Protocol() == ProtocolHTTPS {
		err = s.WebServer.ServeTLS(listener, s.CertFile, s.KeyFile)
	} else {
		err = s.WebServer.Serve(listener)
	}
	if err != nil && err != http.ErrServerClosed {
		return errs.Wrap(err)
	}
	return nil
}

// Shutdown the server gracefully.
func (s *Server) Shutdown() {
	defer s.Logger.Timef("shutdown of %v", s).End()
	gracePeriod := s.ShutdownGracePeriod
	if gracePeriod <= 0 {
		gracePeriod = time.Minute
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(gracePeriod))
	defer cancel()
	if err := s.WebServer.Shutdown(ctx); err != nil {
		s.Logger.Warn(errs.NewfWithCause(err, "Unable to shutdown %s gracefully", s.Protocol()))
	}
}
