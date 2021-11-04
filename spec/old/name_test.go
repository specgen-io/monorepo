package old

import (
	"gotest.tools/assert"
	"testing"
)

func Test_Name_PascalCase(t *testing.T) {
	name := Name{"some_value", nil}
	assert.Equal(t, name.PascalCase(), "SomeValue")
}

func Test_Name_CamelCase(t *testing.T) {
	name := Name{"some_value", nil}
	assert.Equal(t, name.CamelCase(), "someValue")
}

func Test_Name_SnakeCase(t *testing.T) {
	name := Name{"SomeValue", nil}
	assert.Equal(t, name.SnakeCase(), "some_value")
}

func Test_Name_FlatCase(t *testing.T) {
	name := Name{"SomeValue", nil}
	assert.Equal(t, name.FlatCase(), "somevalue")
}
