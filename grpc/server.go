package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"google.golang.org/grpc"
)

type ServerConfig struct {
	Host string `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	Port string `envconfig:"GRPC_PORT" default:"8000"`
	RPS  int    `envconfig:"RPS"`
}

type Server struct {
	*grpc.Server
	listener net.Listener
}

func NewServer(config ServerConfig, ratelimiter ratelimit.Limiter) (*Server, error) {
	addr := net.JoinHostPort(config.Host, config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("initialize tcp listener: %w", err)
	}
	return &Server{
		Server: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				ratelimit.UnaryServerInterceptor(ratelimiter),
			),
			grpc_middleware.WithStreamServerChain(
				ratelimit.StreamServerInterceptor(ratelimiter),
			),
		),
		listener: listener,
	}, nil
}

func (s *Server) ListenAndServe() error {
	return s.Serve(s.listener)
}
