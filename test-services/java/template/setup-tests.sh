#!/bin/bash +x

{{#server.spring}}
export TEST_PARAMETERS_MODE=true
{{/server.spring}}
{{#server.micronaut}}
export TEST_COMMA_SEPARATED_FORM_PARAMS_MODE=true
{{/server.micronaut}}
export TEST_NO_FORM_DATA=false