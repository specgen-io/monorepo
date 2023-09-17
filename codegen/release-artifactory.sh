#!/bin/bash
set -e

if [ -n "$1" ]; then
    VERSION=$1
else
    echo 'Version not set'
    exit 1
fi

if [[ $JFROG_USER == "" ]]; then
    echo '$JFROG_USER variable is not set'
    exit 1
fi

if [[ $JFROG_PASS == "" ]]; then
    echo '$JFROG_PASS variable is not set'
    exit 1
fi

echo "Releasing $VERSION"

RELEASE_NAME=v$VERSION

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
