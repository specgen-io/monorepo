package yamlx

import (
	"gopkg.in/specgen-io/yaml.v3"
)

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

func (self *YamlArray) AddWithComment(value interface{}, comment *string) error {
	valueNode := yaml.Node{}
	err := valueNode.Encode(value)
	if err != nil {
		return err
	}
	if comment != nil {
		valueNode.LineComment = *comment
	}
	self.Node.Content = append(self.Node.Content, &valueNode)
	return nil
}

func (self *YamlArray) Length() int {
	return len(self.Node.Content)
}
