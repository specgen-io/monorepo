#!/bin/bash -x

VERSION="0.0.0"
if [ -n "$1" ]; then
    VERSION=$1
fi

echo "Building version: $VERSION"

platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
for platform in "${platforms[@]}"
do
    echo "Building platform: $platform"

    # parse platforms
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    EXEC_NAME="specgen"
    if [ $GOOS = "windows" ]; then
        EXEC_NAME+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w" -o ./dist/${GOOS}_${GOARCH}/${EXEC_NAME} specgen.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    rm -rf ./output
done

echo "Done building version: $VERSION"
