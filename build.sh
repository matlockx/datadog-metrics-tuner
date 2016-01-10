#!/bin/sh

set -e -o pipefail

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/go/src/app -e DOCKER_IMAGE=matlockx/datadog-metrics-tuner matlockx/golang-builder datadog-metrics-tuner

docker push matlockx/datadog-metrics-tuner
