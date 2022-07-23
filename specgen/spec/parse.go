package spec

import (
	"bytes"
	"gopkg.in/specgen-io/yaml.v3"
)

type SpecParseResult struct {
	Spec     *Spec
	Messages Messages
}

func ReadSpec(data []byte) (*Spec, *Messages, error) {
	return ReadSpecWithOptions(SpecOptionsDefault, data)
}

type SpecOptions struct {
	AddErrors bool
}

var SpecOptionsDefault = SpecOptions{true}

func ReadSpecWithOptions(options SpecOptions, data []byte) (*Spec, *Messages, error) {
	allMessages := NewMessages()
	data, messages, err := checkSpecVersion(data)
	allMessages.AddAll(messages.Items...)
	if err != nil {
		return nil, allMessages, err
	}

	spec, messages, err := unmarshalSpec(data)
	allMessages.AddAll(messages.Items...)
	if err != nil {
		return nil, allMessages, err
	}

	messages, err = enrich(options, spec)
	allMessages.AddAll(messages.Items...)
	if err != nil {
		return nil, allMessages, err
	}

	messages, err = validate(spec)
	allMessages.AddAll(messages.Items...)
	if err != nil {
		return nil, allMessages, err
	}

	return spec, allMessages, nil
}

func WriteSpec(spec *Spec) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(spec)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
