package tests

import (
	"gotest.tools/assert"
	"testing"
	"the-client/spec/v2/echo"
	"the-client/spec/v2/models"
)

func Test_V2_EchoBodyModel(t *testing.T) {
	client := echo.NewClient(serviceUrl)

	expectedMessage := &models.Message{true, "the string"}
	response, err := client.EchoBodyModel(&models.Message{true, "the string"})

	assert.NilError(t, err)
	assert.DeepEqual(t, expectedMessage, response)
}
