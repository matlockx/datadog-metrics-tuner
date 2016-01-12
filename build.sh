#!/bin/sh

set -e -o pipefail

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/go/src/app -e DOCKER_IMAGE=matlockx/datadog-metrics-tuner matlockx/golang-builder -ldflags "-X=main.buildTime=2016-01-12_11:17:04AM -X=main.gitHash=80cb93c613cf1d5db4cc6c44e2ef27c4636a066a" -o datadog-metrics-tuner

docker push matlockx/datadog-metrics-tuner
