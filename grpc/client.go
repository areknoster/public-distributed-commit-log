package grpc

import (
	"net"

	"google.golang.org/grpc"
)

type ConnConfig struct {
	Host string `envconfig:"GRPC_HOST" default:"localhost"`
	Port string `envconfig:"GRPC_PORT" default:"8000"`
}

func Dial(config ConnConfig) (*grpc.ClientConn, error) {
	return grpc.Dial(
		net.JoinHostPort(config.Host, config.Port),
		grpc.WithInsecure(),
	)
}
