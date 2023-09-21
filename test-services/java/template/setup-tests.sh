#!/bin/bash +x

{{#server.spring}}
export TEST_PARAMETERS_MODE=true
{{/server.spring}}
export TEST_NO_FORM_DATA=true