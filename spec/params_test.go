package spec

import (
	"strings"
	"testing"
)

func Test_QueryParams_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
param1: string # some param
param2: int
`, "\n")
	var params QueryParams
	checkUnmarshalMarshal(t, expectedYaml, &params)
}

func Test_HeaderParams_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
Authorization: string # some param
Accept-Language: string # some param
Some-Header: string # some param
`, "\n")
	var params HeaderParams
	checkUnmarshalMarshal(t, expectedYaml, &params)
}
