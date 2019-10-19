#!/bin/bash

BUILDNUM="0"
if [ -n "$1" ]; then
    BUILDNUM=$1
fi
VERSION="0.$BUILDNUM"

echo "Building version: $VERSION"

dep ensure

mkdir -p ./zips

platforms=("windows/amd64" "darwin/amd64" "linux/amd64")
for platform in "${platforms[@]}"
do
    # parse platforms
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    ###### build and package spec
    output_name="specgen_${GOOS}_${GOARCH}"

    exec_name="specgen"
    if [ $GOOS = "windows" ]; then
        exec_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-s -w -X specgen/version.Current=$VERSION" -o $exec_name specgen.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    zip "./zips/$output_name.zip" $exec_name -q

    rm -rf ./output $exec_name "$output_name.zip"
done