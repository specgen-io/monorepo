#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

echo "Releasing $VERSION"

if [ -n "$2" ]; then
    GITHUB_TOKEN=$2
else
    echo 'Github token not set'
    exit 1
fi

zip "./specgen_darwin_amd64.zip" "./dist/darwin_amd64/specgen" -q -j
zip "./specgen_linux_amd64.zip" "./dist/linux_amd64/specgen" -q -j
zip "./specgen_windows_amd64.zip" "./dist/windows_amd64/specgen.exe" -q -j

RELEASE_NAME=v$VERSION

go get github.com/aktau/github-release

echo "Creating release in Github: $RELEASE_NAME"
set +e
$GOPATH/bin/github-release release --security-token $GITHUB_TOKEN --user specgen-io --repo specgen --tag $RELEASE_NAME --target v2
set -e

sleep 10

echo "Releasing specgen_darwin_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo specgen --tag $RELEASE_NAME --name specgen_darwin_amd64.zip  --file specgen_darwin_amd64.zip
echo "Releasing specgen_linux_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo specgen --tag $RELEASE_NAME --name specgen_linux_amd64.zip   --file specgen_linux_amd64.zip
echo "Releasing specgen_windows_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user specgen-io --repo specgen --tag $RELEASE_NAME --name specgen_windows_amd64.zip --file specgen_windows_amd64.zip

echo "Done releasing $RELEASE_NAME"
