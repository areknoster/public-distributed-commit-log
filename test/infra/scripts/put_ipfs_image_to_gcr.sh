#!/bin/bash

set -e


TAG="master-2022-01-14-4403946"
#modify for new image
docker pull ipfs/go-ipfs:${TAG}
docker tag ipfs/go-ipfs:${TAG} eu.gcr.io/pdcl-testing/go-ipfs:${TAG}
docker push  eu.gcr.io/pdcl-testing/go-ipfs:${TAG}
