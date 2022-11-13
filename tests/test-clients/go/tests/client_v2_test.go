package tests

import (
	"gotest.tools/assert"
	"test-client/spec/v2/echo"
	"test-client/spec/v2/models"
	"testing"
)

func Test_V2_EchoBodyModel(t *testing.T) {
	client := echo.NewClient(serviceUrl)

	expectedMessage := &models.Message{true, "the string"}
	response, err := client.EchoBodyModel(&models.Message{true, "the string"})

	assert.NilError(t, err)
	assert.NilError(t, err, response)
	assert.DeepEqual(t, expectedMessage, response)
}
