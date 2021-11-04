package spec

import (
	"strings"
	"testing"
)

func Test_NamedDefinitions_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
prop1: string # some field
prop2: int?
`, "\n")
	var definitions NamedDefinitions
	checkUnmarshalMarshal(t, expectedYaml, &definitions)
}
