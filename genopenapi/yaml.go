package genopenapi

import "gopkg.in/yaml.v2"

type YamlMap struct {
	Yaml []yaml.MapItem
}

func Map() *YamlMap {
	return &YamlMap{make([]yaml.MapItem, 0)}
}

func (self *YamlMap) Set(name string, value interface{}) *YamlMap {
	self.Yaml = append(self.Yaml, yaml.MapItem{name, value})
	return self
}

func (self *YamlMap) String() string {
	data, _ := yaml.Marshal(self.Yaml)
	return string(data)
}

type YamlArray struct {
	Yaml []interface{}
}

func Array() *YamlArray {
	return &YamlArray{make([]interface{}, 0)}
}

func (self *YamlArray) Add(value interface{}) *YamlArray {
	self.Yaml = append(self.Yaml, value)
	return self
}

func (self *YamlArray) Length() int {
	return len(self.Yaml)
}
