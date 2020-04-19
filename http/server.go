package http

import (
	"net"
	"net/http"
)

// Server represents a HTTP server.
type Server struct {
	ln      net.Listener
	Handler *Handler
	Addr    string
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	go func() { http.Serve(s.ln, s.Handler) }()

	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}
