package old

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"strings"
)

func getMappingKey(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(mapping.Content)/2; i++ {
		keyNode := mapping.Content[i*2]
		if keyNode.Value == key {
			return keyNode
		}
	}
	return nil
}

func getMappingValue(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i < len(mapping.Content)/2; i++ {
		keyNode := mapping.Content[i*2]
		if keyNode.Value == key {
			return mapping.Content[i*2+1]
		}
	}
	return nil
}

func decodeStringOptional(node *yaml.Node, key string) (*string, error) {
	strNode := getMappingValue(node, key)
	if strNode == nil {
		return nil, nil
	}
	str := ""
	err := strNode.DecodeWith(decodeStrict, &str)
	if err != nil {
		return nil, yamlError(strNode, "failed to decode string value")
	}
	return &str, err
}

func getDescription(node *yaml.Node) *string {
	if node == nil {
		return nil
	}
	lineComment := node.LineComment
	lineComment = strings.TrimLeft(lineComment, "#")
	lineComment = strings.TrimSpace(lineComment)
	if lineComment == "" {
		return nil
	}
	return &lineComment
}

func yamlError(node *yaml.Node, message string) error {
	return errors.New(fmt.Sprintf("yaml: line %d: %s", node.Line, message))
}

func contains(values []string, node *yaml.Node) bool {
	for _, v := range values {
		if node.Value == v {
			return true
		}
	}
	return false
}

var decodeStrict = yaml.NewDecodeOptions().KnownFields(true)

var decodeLooze = yaml.NewDecodeOptions().KnownFields(false)