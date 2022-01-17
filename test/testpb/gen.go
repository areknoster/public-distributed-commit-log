// Package testpb contains proto definitions used in testing.
package testpb

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. test.proto
