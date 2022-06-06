#!/bin/bash +x

VERSION="0.0.0"
if [ -n "$1" ]; then
    VERSION=$1
fi

echo "Building version: $VERSION"

build()
{
  GOOS=$1
  GOARCH=$2
  EXECNAME=$3

  echo "Building ${GOOS}_${GOARCH}/${EXECNAME}"
  go build -ldflags "-s -w" -o ./dist/${GOOS}_${GOARCH}/${EXECNAME} specgen.go
  if [ $? -ne 0 ]; then
      echo "An error has occurred while building ${GOOS}_${GOARCH}/${EXECNAME}! Aborting the script execution..."
      exit 1
  fi
  echo 'Successfully built'
}

build windows amd64 specgen.exe
build darwin amd64 specgen
build darwin arm64 specgen
build linux amd64 specgen

echo "Done building version: $VERSION"
