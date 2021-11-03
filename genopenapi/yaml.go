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

func Map(pairs ...Pair) *YamlMap {
	self :=  &YamlMap{yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{}}}
	err := self.AddAll(pairs...)
	if err != nil {
		panic(err)
	}
	return self
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

type Pair struct {
	Key   interface{}
	Value interface{}
}

func (yamlMap *YamlMap) Add(key interface{}, value interface{}) error {
	return yamlMap.AddAll(Pair{key, value})
}

func (yamlMap *YamlMap) AddAll(pairs ...Pair) error {
	for _, pair := range pairs {
		keyNode, valueNode, err := encodeKeyValue(pair.Key, pair.Value)
		if err != nil {
			return err
		}
		yamlMap.Node.Content = append(yamlMap.Node.Content, keyNode, valueNode)
	}
	return nil
}

func (self *YamlMap) Length() int {
	return len(self.Node.Content) / 2
}

type YamlArray struct {
	Node yaml.Node
}

func Array(items ...interface{}) *YamlArray {
	self := &YamlArray{yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{}}}
	self.Add(items...)
	return self
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
