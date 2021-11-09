package spec

import (
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type Enum struct {
	Items EnumItems `yaml:"enum"`
}

type EnumItem struct {
	Value       string
	Description *string
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
			if valueNode.Kind != yaml.ScalarNode {
				return yamlError(valueNode, "enum item has to be scalar value")
			}
			array[index] = NamedEnumItem{Name: itemName, EnumItem: EnumItem{valueNode.Value, getDescription(valueNode)}}
		}
		*value = array
	}

	return nil
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
		yamlArray := yamlx.Array()
		for index := 0; index < len(value); index++ {
			item := value[index]
			err := yamlArray.AddWithComment(item.Name, item.Description)
			if err != nil {
				return nil, err
			}
		}
		return yamlArray.Node, nil
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
