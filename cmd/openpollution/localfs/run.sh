#!/bin/bash

trap "killall background" EXIT

GRPC_HOST="localhost" go run ./sentinel &
go run ./random-producer
