package old

import (
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type EnumItem struct {
	Value       string  `yaml:"value"`
	Description *string `yaml:"description"`
}

type NamedEnumItem struct {
	Name Name
	EnumItem
}

type EnumItems []NamedEnumItem

func (value *EnumItems) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode && node.Kind != yaml.MappingNode {
		return yamlError(node, "enum items should be either list or mapping")
	}

	if node.Kind == yaml.SequenceNode {
		count := len(node.Content)
		array := make(EnumItems, count)
		for index := 0; index < count; index++ {
			itemNode := node.Content[index]
			itemName := Name{}
			err := itemNode.DecodeWith(decodeStrict, &itemName)
			if err != nil {
				return err
			}
			err = itemName.Check(SnakeCase)
			if err != nil {
				return err
			}
			array[index] = NamedEnumItem{Name: itemName, EnumItem: EnumItem{Value: itemName.Source, Description: getDescription(itemNode)}}
		}
		*value = array
	}

	if node.Kind == yaml.MappingNode {
		count := len(node.Content) / 2
		array := make(EnumItems, count)
		for index := 0; index < count; index++ {
			keyNode := node.Content[index*2]
			valueNode := node.Content[index*2+1]
			itemName := Name{}
			err := keyNode.DecodeWith(decodeStrict, &itemName)
			if err != nil {
				return err
			}
			err = itemName.Check(SnakeCase)
			if err != nil {
				return err
			}
			item := &EnumItem{}
			if valueNode.Kind == yaml.ScalarNode {
				item.Value = valueNode.Value
				item.Description = getDescription(valueNode)
			} else {
				err = valueNode.DecodeWith(decodeStrict, item)
				if err != nil {
					return err
				}
			}
			if item.Value == "" {
				item.Value = itemName.Source
			}
			array[index] = NamedEnumItem{Name: itemName, EnumItem: *item}
		}
		*value = array
	}

	return nil
}

type Enum struct {
	Items       EnumItems `yaml:"enum"`
	Description *string   `yaml:"description"`
}

func itemsSameAsValues(value EnumItems) bool {
	for _, item := range value {
		if item.Name.Source != item.Value {
			return false
		}
	}
	return true
}

func (value EnumItems) MarshalYAML() (interface{}, error) {
	if itemsSameAsValues(value) {
		node := yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: make([]*yaml.Node, len(value)),
		}
		for index := 0; index < len(value); index++ {
			item := value[index]
			itemNode := yaml.Node{}
			err := itemNode.Encode(item.Name)
			if err != nil {
				return nil, err
			}
			if item.Description != nil {
				itemNode.LineComment = *item.Description
			}
			node.Content[index] = &itemNode
		}
		return node, nil
	} else {
		yamlMap := yamlx.Map()
		for index := 0; index < len(value); index++ {
			item := value[index]
			err := yamlMap.AddWithComment(item.Name, item.Value, item.Description)
			if err != nil {
				return nil, err
			}
		}
		return yamlMap.Node, nil
	}
}

func (value Enum) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Add("enum", value.Items)
	return yamlMap.Node, nil
}
