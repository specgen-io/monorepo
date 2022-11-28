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
	notFoundMessages := []string{}
	for _, expected := range expected {
		found := messages.Contains(func(m Message) bool {
			sameLevel := m.Level == expected.Level
			messageMatched := strings.HasPrefix(m.Message, expected.Message)
			sameLocation := true
			if expected.Location != nil || m.Location != nil {
				if expected.Location != nil && m.Location != nil {
					sameLocation = m.Location.Line == expected.Location.Line && m.Location.Column == expected.Location.Column
				} else {
					sameLocation = false
				}
			}
			return sameLevel && messageMatched && sameLocation
		})
		if !found {
			notFoundMessages = append(notFoundMessages, expected.String())
		}
	}
	if len(notFoundMessages) > 0 {
		actualMessages := []string{}
		for _, msg := range messages.Items {
			actualMessages = append(actualMessages, msg.String())
		}
		t.Errorf("Expected messages were not found:\n%s\nFound:\n%s", strings.Join(notFoundMessages, "\n"), strings.Join(actualMessages, "\n"))
	}
	assert.Equal(t, len(expected), len(messages.Items), `Expected %d messages, found %d`, len(expected), len(messages.Items))
}

const specificationMetaLines = 5

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
