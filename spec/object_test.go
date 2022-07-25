package spec

import (
	"strings"
	"testing"
)

func Test_Object_Marshal(t *testing.T) {
	expectedYaml := strings.TrimLeft(`
object:
  prop1: string
  prop2: string # some field
`, "\n")
	var model Object
	checkUnmarshalMarshal(t, expectedYaml, &model)
}
