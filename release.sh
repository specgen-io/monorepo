#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

echo 'Releasing $VERSION'

if [ -n "$2" ]; then
    GITHUB_TOKEN=$2
else
    echo 'Github token not set'
    exit 1
fi

RELEASE_NAME=$VERSION

go get github.com/aktau/github-release

$GOPATH/bin/github-release release --security-token $GITHUB_TOKEN --user ModaOperandi --repo specgen --tag $RELEASE_NAME

$GOPATH/bin/github-release upload  --security-token $GITHUB_TOKEN --user ModaOperandi --repo specgen --tag $RELEASE_NAME --name specgen_darwin_amd64.zip  --file zips/specgen_darwin_amd64.zip
$GOPATH/bin/github-release upload  --security-token $GITHUB_TOKEN --user ModaOperandi --repo specgen --tag $RELEASE_NAME --name specgen_linux_amd64.zip   --file zips/specgen_linux_amd64.zip
$GOPATH/bin/github-release upload  --security-token $GITHUB_TOKEN --user ModaOperandi --repo specgen --tag $RELEASE_NAME --name specgen_windows_amd64.zip --file zips/specgen_windows_amd64.zip

echo 'Done releasing $VERSION'
