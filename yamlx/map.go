package yamlx

import (
	"gopkg.in/specgen-io/yaml.v3"
	"reflect"
)

type YamlMap struct {
	Node yaml.Node
}

func Map(pairs ...Pair) *YamlMap {
	self := &YamlMap{yaml.Node{Kind: yaml.MappingNode, Content: []*yaml.Node{}}}
	err := self.AddAll(pairs...)
	if err != nil {
		panic(err)
	}
	return self
}

func (self *YamlMap) MarshalYAML() (interface{}, error) {
	return &self.Node, nil
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

func (yamlMap *YamlMap) Add(key interface{}, value interface{}) error {
	return yamlMap.AddAll(Pair{key, value})
}

func (yamlMap *YamlMap) AddRaw(key interface{}, value string) error {
	keyNode := yaml.Node{}
	err := keyNode.Encode(key)
	if err != nil {
		return err
	}
	valueNode := String(value)
	yamlMap.Node.Content = append(yamlMap.Node.Content, &keyNode, &valueNode)
	return nil
}

func (yamlMap *YamlMap) AddOmitNil(key interface{}, value interface{}) error {
	if !reflect.ValueOf(value).IsNil() {
		yamlMap.Add(key, value)
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
