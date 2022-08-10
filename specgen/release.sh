#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

if [ -n "$2" ]; then
    GITHUB_ORG=$2
else
    echo 'Github org is not set'
    exit 1
fi

if [ -n "$3" ]; then
    GITHUB_REPO=$3
else
    echo 'Github repo is not set'
    exit 1
fi

if [ -n "$4" ]; then
    GITHUB_TOKEN=$4
else
    echo 'Github token is not set'
    exit 1
fi

echo "Releasing $VERSION"

zip "./specgen_darwin_amd64.zip" "./dist/darwin_amd64/specgen" -q -j
zip "./specgen_darwin_arm64.zip" "./dist/darwin_arm64/specgen" -q -j
zip "./specgen_linux_amd64.zip" "./dist/linux_amd64/specgen" -q -j
zip "./specgen_windows_amd64.zip" "./dist/windows_amd64/specgen.exe" -q -j

RELEASE_NAME=v$VERSION

go get github.com/aktau/github-release

echo "Creating release in Github: $RELEASE_NAME"
set +e
$GOPATH/bin/github-release release --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME
set -e

sleep 10

echo "Releasing specgen_darwin_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_darwin_amd64.zip  --file specgen_darwin_amd64.zip
echo "Releasing specgen_darwin_arm64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_darwin_arm64.zip  --file specgen_darwin_arm64.zip
echo "Releasing specgen_linux_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_linux_amd64.zip   --file specgen_linux_amd64.zip
echo "Releasing specgen_windows_amd64.zip"
$GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_windows_amd64.zip --file specgen_windows_amd64.zip

echo "Done releasing $RELEASE_NAME"
