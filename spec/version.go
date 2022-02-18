package spec

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/specgen-io/specgen/v2/spec/old"
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
		messages.Add(Warning("warning: unexpected spec format version: %s; please format you spec to format %s", *specVersion, SpecVersion).At(locationFromNode(versionNode)))
		messages.Add(Info("will try to convert spec to format %s on the fly...", SpecVersion))

		convertedData, err := old.FormatSpec(data, SpecVersion)
		if err != nil {
			messages.Add(Error(`format conversion failed from: %s to: %s with error: %s`, SpecVersion, *specVersion, err.Error()).At(locationFromNode(versionNode)))
			return nil, messages, errors.New(fmt.Sprintf(`failed to convert specification to format %s`, SpecVersion))
		}
		messages.Add(Info("convert spec to format %s succeeded", SpecVersion))
		return convertedData, messages, nil
	}
	return data, messages, nil
}
