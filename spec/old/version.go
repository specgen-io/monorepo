package old

import (
	"bytes"
	"errors"
	"fmt"
)
import "gopkg.in/specgen-io/yaml.v3"

var SpecVersion = "2"

func GetSpecVersion(data []byte) (*string, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&node)
	if err != nil {
		return nil, err
	}

	rootNode := node.Content[0]
	idlVersion := getMappingValue(rootNode, "idl_version")
	if idlVersion != nil {
		return &idlVersion.Value, nil
	}
	specVersion := getMappingValue(rootNode, "spec")
	if specVersion != nil {
		return &specVersion.Value, nil
	}
	return nil, errors.New(`can't find old version field, should be either "idl_version" for old v1 or "spec" for later old versions`)
}

func checkSpecVersion(data []byte) ([]byte, error) {
	specVersion, err := GetSpecVersion(data)
	if err != nil {
		return nil, err
	}

	if *specVersion == "0" || *specVersion == "1" {
		var node yaml.Node
		decoder := yaml.NewDecoder(bytes.NewReader(data))
		err = decoder.Decode(&node)
		if err != nil {
			return nil, err
		}

		rootNode := node.Content[0]
		idlVersionKey := getMappingKey(rootNode, "idl_version")
		idlVersionValue := getMappingValue(rootNode, "idl_version")
		idlVersionKey.Value = "spec"
		idlVersionValue.Value = "2"
		operations := getMappingKey(rootNode, "operations")
		if operations != nil {
			operations.Value = "http"
		}
		serviceName := getMappingKey(rootNode, "service_name")
		serviceName.Value = "name"
		data, err = yaml.Marshal(&node)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else if *specVersion != SpecVersion {
		return nil, fmt.Errorf("unexpected old version, expected: %s, found: %s", SpecVersion, *specVersion)
	}
	return data, nil
}
