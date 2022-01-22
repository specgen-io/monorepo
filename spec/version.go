package spec

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/specgen-io/specgen/v2/console"
	"github.com/specgen-io/specgen/v2/spec/old"
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
	if specVersion != nil {
		return specVersion, nil
	}

	specVersion, err = decodeStringOptional(rootNode, "idl_version")
	if err != nil {
		return nil, err
	}
	if specVersion != nil {
		return specVersion, nil
	}

	return nil, errors.New(`can't find "spec" or legacy "idl_version" containing spec version`)
}

func checkSpecVersion(data []byte) ([]byte, error) {
	specVersion, err := getSpecVersion(data)
	if err != nil {
		return nil, err
	}
	if *specVersion != SpecVersion {
		console.PrintLnF("warning: unexpected spec format version: %s; please format you spec to format %s", *specVersion, SpecVersion)
		console.PrintLnF("trying to convert spec to format %s on the fly...", SpecVersion)
		convertedData, err := old.FormatSpec(data, SpecVersion)
		if err != nil {
			return nil, errors.New(fmt.Sprintf(`unexpected spec version, expected: %s, found: %s, conversion failed: %s`, SpecVersion, *specVersion, err.Error()))
		}
		console.PrintLnF("convert spec to format %s succeeded", SpecVersion)
		return convertedData, nil
	}
	return data, nil
}
