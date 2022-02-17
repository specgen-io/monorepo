package spec

import (
	"bytes"
	"errors"
	"gopkg.in/specgen-io/yaml.v3"
	"strconv"
	"strings"
)

type SpecParseResult struct {
	Spec     *Spec
	Messages Messages
}

func ReadSpec(data []byte) (*Spec, Messages, error) {
	data, err := checkSpecVersion(data)
	if err != nil {
		message := convertError(err)
		if message != nil {
			return nil, Messages{*message}, errors.New("failed to read specification")
		}
		return nil, Messages{}, err
	}

	spec, err := unmarshalSpec(data)
	if err != nil {
		message := convertError(err)
		if message != nil {
			return nil, Messages{*message}, errors.New("failed to read specification")
		}
		return nil, Messages{}, err
	}

	messages, err := enrich(spec)
	if err != nil {
		return nil, messages, err
	}

	messages, err = validate(spec)
	if err != nil {
		return nil, messages, err
	}

	return spec, messages, nil
}

func convertError(err error) *Message {
	if strings.HasPrefix(err.Error(), "yaml: line ") {
		parts := strings.SplitN(strings.TrimPrefix(err.Error(), "yaml: line "), ":", 2)
		line, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil
		}
		return &Message{LevelError, strings.TrimSpace(parts[1]), &Location{line, 0}}
	}
	return nil
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

func ReadMeta(data []byte) (*Meta, error) {
	data, err := checkSpecVersion(data)
	if err != nil {
		return nil, err
	}
	var meta Meta
	if err := yaml.UnmarshalWith(decodeLooze, data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
