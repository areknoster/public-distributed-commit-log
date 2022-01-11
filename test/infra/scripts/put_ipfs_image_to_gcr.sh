#!/bin/bash

set -e


TAG="master-2022-01-06-2a871ef"
#modify for new image
docker pull ipfs/go-ipfs:${TAG}
docker tag ipfs/go-ipfs:${TAG} eu.gcr.io/pdcl-testing/go-ipfs:${TAG}
docker push  eu.gcr.io/pdcl-testing/go-ipfs:${TAG}
