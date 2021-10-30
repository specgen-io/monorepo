package spec

import (
	"errors"
	"fmt"
	"gopkg.in/specgen-io/yaml.v3"
	"reflect"
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

type YamlMap struct {
	Node yaml.Node
}

func NewYamlMap() *YamlMap {
	return &YamlMap{yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{}}}
}

func encodeKeyValue(key interface{}, value interface{}) (*yaml.Node, *yaml.Node, error) {
	keyNode := yaml.Node{}
	err := keyNode.Encode(key)
	if err != nil {
		return nil, nil, err
	}
	valueNode := yaml.Node{}
	err = valueNode.Encode(value)
	if err != nil {
		return nil, nil, err
	}
	return &keyNode, &valueNode, nil
}

func (yamlMap *YamlMap) Add(key interface{}, value interface{}) error {
	keyNode, valueNode, err := encodeKeyValue(key, value)
	if err != nil {
		return err
	}
	yamlMap.Node.Content = append(yamlMap.Node.Content, keyNode, valueNode)
	return nil
}

func (yamlMap *YamlMap) AddOmitNil(key interface{}, value interface{}) error {
	if !reflect.ValueOf(value).IsNil() {
		keyNode, valueNode, err := encodeKeyValue(key, value)
		if err != nil {
			return err
		}
		yamlMap.Node.Content = append(yamlMap.Node.Content, keyNode, valueNode)
	}
	return nil
}

func (yamlMap *YamlMap) AddWithComment(key interface{}, value interface{}, comment *string) error {
	keyNode, valueNode, err := encodeKeyValue(key, value)
	if err != nil {
		return err
	}
	if comment != nil {
		if valueNode.Kind == yaml.ScalarNode {
			valueNode.LineComment = *comment
		} else {
			keyNode.LineComment = *comment
		}
	}
	yamlMap.Node.Content = append(yamlMap.Node.Content, keyNode, valueNode)
	return nil
}

func (yamlMap *YamlMap) Merge(value interface{}) error {
	valueNode := yaml.Node{}
	err := valueNode.Encode(value)
	if err != nil {
		return err
	}
	for index := 0; index < len(valueNode.Content); index++ {
		yamlMap.Node.Content = append(yamlMap.Node.Content, valueNode.Content[index])
	}

	return nil
}

func YamlString(value string) yaml.Node {
	return yaml.Node{Kind: yaml.ScalarNode, Value: value}
}
