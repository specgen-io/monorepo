package spec

import (
	"bytes"
	"errors"
	"fmt"
)
import "gopkg.in/specgen-io/yaml.v3"

var SpecVersion = "2.1"

func getSpecVersion(data []byte) (*string, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&node)
	if err != nil {
		return nil, err
	}

	rootNode := node.Content[0]
	specVersion, err := decodeStringOptional(rootNode, "spec")
	if err != nil {
		return nil, err
	}

	if specVersion == nil {
		return nil, errors.New(`Can't find "spec" field with spec version`)
	}
	return specVersion, nil
}

func checkSpecVersion(data []byte) error {
	specVersion, err := getSpecVersion(data)
	if err != nil {
		return err
	}
	if *specVersion != SpecVersion {
		return fmt.Errorf("unexpected spec version, expected: %s, found: %s", SpecVersion, *specVersion)
	}
	return nil
}
