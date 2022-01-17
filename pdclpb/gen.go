// Package pdclpb contains proto definitions of PDCL structures.
package pdclpb

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. pdcl.proto
