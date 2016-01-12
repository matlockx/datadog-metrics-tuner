#!/bin/sh

set -e -o pipefail

if [ -n "$(git status -s)" ]; then
	echo "Your working directory isn't clean. Commit or stash your changes and continue."
	exit 1
fi

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $(pwd):/go/src/app -e DOCKER_IMAGE=matlockx/datadog-metrics-tuner matlockx/golang-builder -ldflags "-X=main.buildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X=main.gitHash=$(git rev-parse HEAD)" -o datadog-metrics-tuner

docker push matlockx/datadog-metrics-tuner
