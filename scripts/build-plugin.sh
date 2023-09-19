#!/bin/bash +x

VERSION="0.0.0"
if [ -n "$1" ]; then
    VERSION=$1
fi

NAME="specgen"

echo "Building $NAME version: $VERSION"

echo "Stamping version to the source code: $VERSION"
cat <<END > ./version/version.go
package version

var Current = "$VERSION"

END

build()
{
  GOOS=$1
  GOARCH=$2

  EXECNAME=${NAME}
  if [[ $GOOS == windows ]]; then
    EXECNAME=${NAME}.exe
  fi

  echo "Building ${GOOS}_${GOARCH}/${EXECNAME}"
  env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o ./dist/${GOOS}_${GOARCH}/${EXECNAME} main.go
  if [ $? -ne 0 ]; then
      echo "An error has occurred while building ${GOOS}_${GOARCH}/${EXECNAME}! Aborting the script execution..."
      exit 1
  fi
  echo 'Successfully built'
}

build windows amd64
build darwin amd64
build darwin arm64
build linux amd64

echo "Done building version: $VERSION"
