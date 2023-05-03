package server

import (
	"context"
	stdErr "errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/thanhfphan/eventstore/pkg/errors"
)

// Server provides a graceful shutdown
type Server struct {
	ip       string
	port     string
	listener net.Listener
}

// New create a new server listening on the provided port. It will starts the listener but
// does not start the server. If an empty port is given, the server randomly chooses one.
func New(port string) (*Server, error) {
	addr := fmt.Sprintf(":" + port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		ip:       listener.Addr().(*net.TCPAddr).IP.String(),
		port:     strconv.Itoa(listener.Addr().(*net.TCPAddr).Port),
		listener: listener,
	}, nil
}

func (s *Server) ServeHTTP(ctx context.Context, srv *http.Server) error {
	errCh := make(chan error, 1)
	go func() {
		<-ctx.Done()

		shutdownCtx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		errCh <- srv.Shutdown(shutdownCtx)
	}()

	err := srv.Serve(s.listener)
	if err != nil && !stdErr.Is(err, http.ErrServerClosed) {
		return errors.New("failed to serve: %w", err)
	}

	err = <-errCh
	return err
}

func (s *Server) ServeHTTPHandler(ctx context.Context, handler http.Handler) error {
	return s.ServeHTTP(ctx, &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	})
}

func (s *Server) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

func (s *Server) IP() string {
	return s.ip
}

func (s *Server) Port() string {
	return s.port
}
