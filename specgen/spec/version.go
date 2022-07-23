package spec

import (
	"bytes"
	"errors"
	"fmt"
)
import "gopkg.in/specgen-io/yaml.v3"

var SpecVersion = "2.1"

func getSpecVersion(data []byte) (*string, *yaml.Node, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&node)
	if err != nil {
		return nil, nil, err
	}

	rootNode := node.Content[0]

	specVersion, specNode, err := decodeStringNodeOptional(rootNode, "spec")
	if err != nil {
		return nil, specNode, err
	}
	if specVersion != nil {
		return specVersion, specNode, nil
	}

	specVersion, specNode, err = decodeStringNodeOptional(rootNode, "idl_version")
	if err != nil {
		return nil, specNode, err
	}
	if specVersion != nil {
		return specVersion, specNode, nil
	}

	return nil, nil, errors.New(`can't find "spec" or legacy "idl_version" containing spec version`)
}

func checkSpecVersion(data []byte) ([]byte, *Messages, error) {
	messages := NewMessages()
	specVersion, versionNode, err := getSpecVersion(data)
	if err != nil {
		messages.Add(convertYamlError(err, versionNode))
		return nil, messages, errors.New("failed to read specification")
	}

	if *specVersion != SpecVersion {
		messages.Add(Error("unexpected spec format version: %s; please format you spec to format %s", *specVersion, SpecVersion).At(locationFromNode(versionNode)))
		return nil, messages, errors.New(fmt.Sprintf(`unexpected spec format version: %s`, *specVersion))
	}
	return data, messages, nil
}
