package spec

import (
	"gotest.tools/assert"
	"strings"
	"testing"
)

type ReadSpecificationCase struct {
	name          string
	specification string
	err           error
	messages      []Message
	check         func(t *testing.T, spec *Spec)
}

func assertMessages(t *testing.T, expected []Message, messages *Messages) {
	assert.Equal(t, len(expected), len(messages.Items), `Expected %s messages, found %s`, len(expected), len(messages.Items))
	for _, expected := range expected {
		found := messages.Contains(func(m Message) bool { return m.Level == expected.Level && m.Message == expected.Message })
		if !found {
			actualMessages := []string{}
			for _, msg := range messages.Items {
				actualMessages = append(actualMessages, msg.String())
			}
			t.Errorf("Expected %s message was not found: '%s'\nFound:\n%s", expected.Level, expected.Message, strings.Join(actualMessages, "\n"))
		}
	}
}

func runReadSpecificationCases(t *testing.T, cases []ReadSpecificationCase) {
	for _, testcase := range cases {
		t.Logf(`Running test case: %s`, testcase.name)

		specificationMeta := `
spec: 2.1
name: testing
version: 1
`

		specification := specificationMeta + testcase.specification
		spec, messages, err := ReadSpec([]byte(specification))

		if testcase.err == nil {
			assert.Equal(t, err, nil)
		} else {
			assert.Error(t, err, testcase.err.Error())
		}
		assertMessages(t, testcase.messages, messages)

		if testcase.check != nil {
			testcase.check(t, spec)
		}
	}
}
