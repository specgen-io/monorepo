package tests

import (
	"encoding/json"
	"gotest.tools/v3/assert"
	"reflect"
	"testing"
	"the-models/spec/v2/models"
)

func TestMessageV2(t *testing.T) {
	data := models.Message{
		"the string",
	}

	jsonStr := `{"field":"the string"}`

	actualJson, err := json.Marshal(data)
	assert.NilError(t, err)
	assert.Equal(t, jsonStr, string(actualJson))

	var actualData models.Message
	err = json.Unmarshal([]byte(jsonStr), &actualData)
	assert.NilError(t, err)
	assert.Equal(t, reflect.DeepEqual(data, actualData), true)
}
