// Package httpshutdown illustrates a possible way to do graceful
// shutdown with net/http. This code is untested.
package httpshutdown

import (
	"net"
	"net/http"
	"sync"
)

// Serve wraps the net/http Server and performs graceful shutdown.
type Server struct {
	Server *http.Server

	w sync.WaitGroup
}

// Serve calls Serve on the underlying http Server,
// but wraps l in another net.Listener that synchronises
// open connections with Wait.
func (s *Server) Serve(l net.Listener) error {
	return s.Server.Serve(&listener{l, &s.w})
}

// Wait waits for all open connections in s to close.
func (s *Server) Wait() {
	s.w.Wait()
}
