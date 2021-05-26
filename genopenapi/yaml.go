package genopenapi

import (
	"bytes"
	yaml "gopkg.in/vsapronov/yaml.v3"
)

type MapItem struct {
	key   string
	value interface{}
}

type YamlMap struct {
	Items []MapItem
}

func (self *YamlMap) MarshalYAML() (interface{}, error) {
	var nodes = []*yaml.Node{}

	for _, item := range self.Items {
		var keyNode yaml.Node
		err := keyNode.Encode(item.key)
		if err != nil {
			return nil, err
		}
		var valueNode yaml.Node
		err = valueNode.Encode(item.value)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &keyNode, &valueNode)
	}

	node := yaml.Node {
		Kind: yaml.MappingNode,
		Tag:  "!!map",
		Content: nodes,
	}
	return &node, nil
}

func Map() *YamlMap {
	return &YamlMap{make([]MapItem, 0)}
}

func (self *YamlMap) Set(name string, value interface{}) *YamlMap {
	self.Items = append(self.Items, MapItem{name, value})
	return self
}

func (self *YamlMap) Length() int {
	return len(self.Items)/2
}

type YamlArray struct {
	Items []interface{}
}

func Array() *YamlArray {
	return &YamlArray{make([]interface{}, 0)}
}

func (self *YamlArray) MarshalYAML() (interface{}, error) {
	return self.Items, nil
}

func (self *YamlArray) Add(value interface{}) *YamlArray {
	self.Items = append(self.Items, value)
	return self
}

func (self *YamlArray) Length() int {
	return len(self.Items)
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