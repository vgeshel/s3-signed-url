#!/bin/sh

set -e

docker run -it --rm -v $(pwd):/work -w /work golang \
       sh -c "go get -d . && CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' s3-signed-url.go"
