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
    OUT_FOLDER=$5
fi

if [ "$GENERATOR" == "models" ]; then
    TESTS_FOLDER=test-${GENERATOR}
else
    TESTS_FOLDER=test-${GENERATOR}s
fi

MAIN_TEMPLATE=./templates/${GENERATOR}-${LANGUAGE}
TEST_TEMPLATE=./${TESTS_FOLDER}/${LANGUAGE}/template
PARAMS_FILE=./${TESTS_FOLDER}/${LANGUAGE}/${PARAMS}.yaml
SPEC_FILE=./${TESTS_FOLDER}/spec.yaml

COMMAND="${RENDR_PATH}rendr file:///${MAIN_TEMPLATE} --root file:///${TEST_TEMPLATE} --noinput --values ${PARAMS_FILE} --set versions.specgen=${SPECGEN_VERSION} --out ${OUT_FOLDER}"
echo "Executing: ${COMMAND}"
$COMMAND

COMMAND="cp ${SPEC_FILE} ${OUT_FOLDER}/spec.yaml"
echo "Executing: ${COMMAND}"
$COMMAND
