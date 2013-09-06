// Package httpshutdown illustrates a possible way to do graceful
// shutdown with net/http. This code is untested.
package httpshutdown

import (
	"net"
	"net/http"
)

// Serve wraps the net/http Server and performs graceful shutdown.
type Server struct {
	server   *http.Server
	listener listener
	stop     chan int
}

func NewServer(s *http.Server, l net.Listener) *Server {
	return &Server{s, listener{Listener: l}, make(chan int, 1)}
}

// Serve calls Serve on the underlying http Server.
// As with package http, Serve returns when the net.Listener in s
// returns an error, but Serve also waits for all open connections
// to close if and only if Shutdown was called.
func (s *Server) Serve() error {
	err := s.server.Serve(&s.listener)
	select {
	case <-s.stop:
		s.listener.wait()
	}
	return err
}

// Shutdown performs a graceful shutdown of s.
// It calls Close on the net.Listener in s.
// Any outstanding requests will complete normally;
// once all open connections have closed, method Serve
// will return.
func (s *Server) Shutdown() {
	select {
	case s.stop <- 1:
	default:
	}
	s.listener.Close()
}
