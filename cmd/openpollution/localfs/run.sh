#!/bin/bash

trap "kill -SIGKILL $(jobs -p)" EXIT

set -e
GRPC_HOST="localhost" go run ./sentinel &
go run ./random-producer &
go run ./print-consumer
