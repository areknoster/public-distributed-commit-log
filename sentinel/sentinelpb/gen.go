// Package sentinelpb contains proto definitions of sentinel services.
package sentinelpb

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. sentinel.proto
