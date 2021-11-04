package old

import (
	"gotest.tools/assert"
	"testing"
)

func Test_Format_Integer_Pass(t *testing.T) {
	err := Integer.Check("123")
	assert.Equal(t, err == nil, true)

	err = Integer.Check("+123")
	assert.Equal(t, err == nil, true)

	err = Integer.Check("-123")
	assert.Equal(t, err == nil, true)

	err = Integer.Check("0")
	assert.Equal(t, err == nil, true)
}

func Test_Format_Integer_Fail(t *testing.T) {
	err := Integer.Check("123.4")
	assert.Equal(t, err != nil, true)

	err = Integer.Check("abcd")
	assert.Equal(t, err != nil, true)

	err = Integer.Check("-")
	assert.Equal(t, err != nil, true)

	err = Integer.Check("+")
	assert.Equal(t, err != nil, true)
}

func Test_Format_JsonFields_Pass(t *testing.T) {
	err := JsonField.Check("name_123")
	assert.Equal(t, err == nil, true)

	err = JsonField.Check("Name123")
	assert.Equal(t, err == nil, true)

	err = JsonField.Check("nameOneTwoThree")
	assert.Equal(t, err == nil, true)
}

func Test_Format_JsonFields_Fail(t *testing.T) {
	err := JsonField.Check("123name")
	assert.Equal(t, err != nil, true)

	err = JsonField.Check("name.123")
	assert.Equal(t, err != nil, true)

	err = JsonField.Check("Name()123")
	assert.Equal(t, err != nil, true)

	err = JsonField.Check("name&One&Two&Three")
	assert.Equal(t, err != nil, true)

	err = JsonField.Check("name$one$two$three")
	assert.Equal(t, err != nil, true)

	err = JsonField.Check("name-one-two-three")
	assert.Equal(t, err != nil, true)
}
