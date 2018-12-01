#!/bin/bash

PROJECT="rcmd"
VERSION="v0.2.0"
BUILD=`date +%Y%m%d`
KEY="ZhRMghgWs&sJrmWs@Sl-CDCN"
function build() {
    os=$1
    arch=$2
    package=${PROJECT}-$3

    echo "build ${package} ..."
    mkdir -p "../releases/${package}"
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -o "../releases/${package}/${PROJECT}" -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X github.com/lcl101/rcmd/core.StrKey=${KEY}" main.go
    # cp ./al.conf "../releases/${package}/al.conf"
    cd ../releases/
    zip -r "./${package}.zip" "./${package}"
    echo "clean ${package}"
    rm -rf "./${package}"
    cd ../al
}



# Linux
build linux amd64 linux
# build linux 386 linux
# build linux arm linux

# OS X Mac
build darwin amd64 macOS

# windows
build windows amd64 win
