#!/bin/bash -x

VERSION="0.0.0"
if [ -n "$1" ]; then
    VERSION=$1
fi

echo "Building version: $VERSION"

mkdir -p ./dist

platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
for platform in "${platforms[@]}"
do
    echo "Building platform: $platform"

    # parse platforms
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    exec_name="specgen"
    if [ $GOOS = "windows" ]; then
        exec_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X github.com/specgen-io/specgen/v2/version.Current=$VERSION" -o $exec_name specgen.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    mkdir -p ./dist/${GOOS}_${GOARCH}
    cp $exec_name ./dist/${GOOS}_${GOARCH}/$exec_name

    rm -rf ./output $exec_name
done

echo "Done building version: $VERSION"
