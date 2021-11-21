#!/bin/bash

trap "kill -SIGKILL $(jobs -p)" EXIT

GRPC_HOST="localhost" go run ./sentinel &
go run ./random-producer
