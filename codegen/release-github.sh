#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

if [[ $GH_TOKEN == "" ]]; then
    echo '$GH_TOKEN variable is not set'
    exit 1
fi

echo "Releasing $VERSION"

RELEASE_NAME=v$VERSION

GH_ORG="specgen-io"
GH_REPO="specgen"

go install github.com/aktau/github-release@v0.10.0

echo "Creating release in Github: $RELEASE_NAME"
set +e
GH_PARAMS="--security-token $GH_TOKEN --user $GH_ORG --repo $GH_REPO --tag $RELEASE_NAME"
github-release release $GH_PARAMS
set -e

sleep 10

GH_RELEASE_URL="https://github.com/${GH_ORG}/${GH_REPO}/releases/tag/${RELEASE_NAME}"
echo "Releasing to Github: ${GH_RELEASE_URL}"

github-release upload --replace $GH_PARAMS --name specgen_darwin_amd64.zip  --file ./specgen_darwin_amd64.zip
github-release upload --replace $GH_PARAMS --name specgen_darwin_arm64.zip  --file ./specgen_darwin_arm64.zip
github-release upload --replace $GH_PARAMS --name specgen_linux_amd64.zip   --file ./specgen_linux_amd64.zip
github-release upload --replace $GH_PARAMS --name specgen_windows_amd64.zip --file ./specgen_windows_amd64.zip

echo "Done releasing to Github: ${GH_RELEASE_URL}"
