#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

if [ -n "$2" ]; then
    TARGET=$2
else
    echo 'Target not set'
    exit 1
fi

echo "Zipping binaries"

zip "./specgen_darwin_amd64.zip" "./dist/darwin_amd64/specgen" -q -j
zip "./specgen_darwin_arm64.zip" "./dist/darwin_arm64/specgen" -q -j
zip "./specgen_linux_amd64.zip" "./dist/linux_amd64/specgen" -q -j
zip "./specgen_windows_amd64.zip" "./dist/windows_amd64/specgen.exe" -q -j

echo "Releasing $VERSION"

RELEASE_NAME=v$VERSION

if [[ $TARGET == *"github"* ]]; then

    GITHUB_ORG="specgen-io"
    GITHUB_REPO="specgen"

    if [[ $GITHUB_TOKEN == "" ]]; then
        echo 'GITHUB_TOKEN variable is not set'
        exit 1
    fi

    go get github.com/aktau/github-release

    echo "Creating release in Github: $RELEASE_NAME"
    set +e
    $GOPATH/bin/github-release release --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME
    set -e

    sleep 10

    echo "Releasing specgen_darwin_amd64.zip"
    $GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_darwin_amd64.zip  --file ./zips/specgen_darwin_amd64.zip
    echo "Releasing specgen_darwin_arm64.zip"
    $GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_darwin_arm64.zip  --file ./zips/specgen_darwin_arm64.zip
    echo "Releasing specgen_linux_amd64.zip"
    $GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_linux_amd64.zip   --file ./zips/specgen_linux_amd64.zip
    echo "Releasing specgen_windows_amd64.zip"
    $GOPATH/bin/github-release upload --replace --security-token $GITHUB_TOKEN --user $GITHUB_ORG --repo $GITHUB_REPO --tag $RELEASE_NAME --name specgen_windows_amd64.zip --file ./zips/specgen_windows_amd64.zip

    echo "Done releasing to Github $RELEASE_NAME"

fi

if [[ $TARGET == *"artifactory"* ]]; then

    ARTFACTORY_URL="https://specgen.jfrog.io/artifactory/binaries/specgen"

    echo "Releasing to Artifactory: $ARTFACTORY_URL/latest"

    curl -u$JFROG_USER:$JFROG_PASS -T specgen_darwin_amd64.zip "$ARTFACTORY_URL/latest/specgen_darwin_amd64.zip"
    curl -u$JFROG_USER:$JFROG_PASS -T specgen_darwin_arm64.zip "$ARTFACTORY_URL/latest/specgen_darwin_arm64.zip"
    curl -u$JFROG_USER:$JFROG_PASS -T specgen_linux_amd64.zip "$ARTFACTORY_URL/latest/specgen_linux_amd64.zip"
    curl -u$JFROG_USER:$JFROG_PASS -T specgen_windows_amd64.zip "$ARTFACTORY_URL/latest/specgen_windows_amd64.zip"

    echo "Releasing to Artifactory: $ARTFACTORY_URL/$RELEASE_NAME"

    curl -u$JFROG_USER:$JFROG_PASS -T specgen_darwin_amd64.zip "$ARTFACTORY_URL/$RELEASE_NAME/specgen_darwin_amd64.zip"
    curl -u$JFROG_USER:$JFROG_PASS -T specgen_darwin_arm64.zip "$ARTFACTORY_URL/$RELEASE_NAME/specgen_darwin_arm64.zip"
    curl -u$JFROG_USER:$JFROG_PASS -T specgen_linux_amd64.zip "$ARTFACTORY_URL/$RELEASE_NAME/specgen_linux_amd64.zip"
    curl -u$JFROG_USER:$JFROG_PASS -T specgen_windows_amd64.zip "$ARTFACTORY_URL/$RELEASE_NAME/specgen_windows_amd64.zip"

    echo "Done releasing to Artifactory $RELEASE_NAME"

fi
