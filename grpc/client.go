package grpc

import (
	"google.golang.org/grpc"
	"net"
)

type ConnConfig struct{
	Host, Port string
}

func Dial(config ConnConfig) (*grpc.ClientConn, error) {
	return grpc.Dial(
		net.JoinHostPort(config.Host, config.Port),
		grpc.WithInsecure(),
	)
}
