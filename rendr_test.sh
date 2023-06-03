#!/bin/bash +x

USAGE="usage: ./rendr_test.sh <GENERATOR> <LANGUAGE> <PARAMS> <SPECGEN_VERSION> [<OUT_FOLDER>]"

if [ -z "$4" ];
then
    echo "Error: not enough params."
    echo "Usage: ./rendr_test.sh <GENERATOR> <LANGUAGE> <PARAMS> <SPECGEN_VERSION> [<OUT_FOLDER>]"
    echo "Example: ./rendr_test.sh service kotlin spring-jackson 0.0.0"
    exit 1
fi

GENERATOR=$1
LANGUAGE=$2
PARAMS=$3
SPECGEN_VERSION=$4
OUT_FOLDER=./the-${GENERATOR}-${LANGUAGE}

if [ -n "$5" ]; then
    OUT_FOLDER=$4
fi

TESTS_FOLDER=test-${GENERATOR}s

MAIN_TEMPLATE=./templates/${GENERATOR}-${LANGUAGE}
TEST_TEMPLATE=./${TESTS_FOLDER}/${LANGUAGE}/template
PARAMS_FILE=./${TESTS_FOLDER}/${LANGUAGE}/${PARAMS}.json
SPEC_FILE=./${TESTS_FOLDER}/spec.yaml

echo "MAIN_TEMPLATE: ${MAIN_TEMPLATE}"
echo "TEST_TEMPLATE: ${TEST_TEMPLATE}"
echo "PARAMS_FILE:   ${PARAMS_FILE}"
echo "SPEC_FILE:     ${SPEC_FILE}"

./rendr file:///${MAIN_TEMPLATE} --root file:///${TEST_TEMPLATE} --noinput --values ${PARAMS_FILE} --set versions.specgen=${SPECGEN_VERSION} --out ${OUT_FOLDER}
cp ${SPEC_FILE} ${OUT_FOLDER}/spec.yaml
