#!/usr/bin/env bash

set -e

DOCKERFILE="$(cat << EOF
FROM golang:latest as build
SHELL ["/bin/bash", "-c"]
WORKDIR /build
ADD . /build
RUN go mod download
RUN cd cmd/pbin && \
    GOOS=darwin \
    GOARCH=amd64 \
    go build \
        -a \
        -o /build/pbin-darwin \
        .
RUN cd cmd/pbin && \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
        -a \
        -o /build/pbin-linux \
        .
RUN cd cmd/pbin && \
    GOOS=windows \
    GOARCH=amd64 \
    go build \
        -a \
        -o /build/pbin-windows \
        .
RUN tar -cvzf build.tgz \
    pbin-linux \
    pbin-darwin \
    pbin-windows
FROM golang:latest
COPY --from=build /build/build.tgz /opt/
CMD ["cat", "/opt/build.tgz"]
EOF
)"

function build
{
    local SCRIPT_PATH="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"
    pushd "${SCRIPT_PATH}"
    while [[ "${PWD}" != '/' && ! -f "go.mod" ]] ; do
        cd ..
    done
    docker build \
        -t build:tmp \
        -f - \
        . <<< "${DOCKERFILE}"
    docker run --rm build:tmp > "${SCRIPT_PATH}/build.tgz"
    popd
}

build

