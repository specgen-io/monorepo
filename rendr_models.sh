#!/bin/bash +x

USAGE="usage: ./rendr_tests.sh <GENERATOR> <PARAMS> <SPECGEN_VERSION>"

if [ -n "$1" ]; then
    GENERATOR=$1
else
    echo "Not enough params, ${USAGE}"
    exit 1
fi

if [ -n "$2" ]; then
    PARAMS=$2
else
    echo "Not enough params, ${USAGE}"
    exit 1
fi

if [ -n "$3" ]; then
    SPECGEN_VERSION=$3
else
    echo "Not enough params, ${USAGE}"
    exit 1
fi

OUT_FOLDER=./the-models-${GENERATOR}
if [ -n "$4" ]; then
    OUT_FOLDER=$4
fi

./rendr file:///./templates/models-${GENERATOR} --root file:///./test-models/${GENERATOR}/template --noinput --values ./test-models/${GENERATOR}/${PARAMS}.json --set versions.specgen=${SPECGEN_VERSION} --out ${OUT_FOLDER}
cp ./test-models/spec.yaml ${OUT_FOLDER}/spec.yaml
