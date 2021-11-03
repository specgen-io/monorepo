package genopenapi

import (
	"bytes"
	yaml "gopkg.in/specgen-io/yaml.v3"
)

type YamlMap struct {
	Node yaml.Node
}

func (self *YamlMap) MarshalYAML() (interface{}, error) {
	return &self.Node, nil
}

func Map() *YamlMap {
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

func (self *YamlMap) Length() int {
	return len(self.Node.Content)/2
}

type YamlArray struct {
	Node yaml.Node
}

func Array() *YamlArray {
	return &YamlArray{yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{}}}
}

func (self *YamlArray) MarshalYAML() (interface{}, error) {
	return self.Node, nil
}

func (self *YamlArray) Add(values ...interface{}) error {
	for _, value := range values {
		valueNode := yaml.Node{}
		err := valueNode.Encode(value)
		if err != nil {
			return err
		}
		self.Node.Content = append(self.Node.Content, &valueNode)
	}
	return nil
}

func (self *YamlArray) Length() int {
	return len(self.Node.Content)
}

func ToYamlString(data interface{}) (string, error) {
	writer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)
	err := encoder.Encode(data)
	if err != nil {
		return "", nil
	}
	return writer.String(), nil
}