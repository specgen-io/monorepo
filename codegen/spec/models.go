package spec

import (
	"errors"
	"gopkg.in/specgen-io/yaml.v3"
	"yamlx"
)

type Model struct {
	Object      *Object
	Enum        *Enum
	OneOf       *OneOf
	Description *string
	Location    *yaml.Node
}

type NamedModel struct {
	Name Name
	Model
	InVersion    *Version
	InHttpErrors *HttpErrors
}

type Models []NamedModel

func (self *Model) IsObject() bool {
	return self.Object != nil && self.Enum == nil && self.OneOf == nil
}

func (self *Model) IsEnum() bool {
	return self.Object == nil && self.Enum != nil && self.OneOf == nil
}

func (self *Model) IsOneOf() bool {
	return self.Object == nil && self.Enum == nil && self.OneOf != nil
}

func (value *Model) UnmarshalYAML(node *yaml.Node) error {
	model := Model{Location: node}
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be mapping")
	}

	description, err := decodeStringOptional(node, "description")
	if err != nil {
		return err
	}
	model.Description = description

	if getMappingKey(node, "enum") != nil {
		enum := Enum{}
		err := node.DecodeWith(decodeLooze, &enum)
		if err != nil {
			return err
		}
		model.Enum = &enum
	} else if getMappingKey(node, "oneOf") != nil {
		oneOf := OneOf{}
		err := node.DecodeWith(decodeLooze, &oneOf)
		if err != nil {
			return err
		}
		model.OneOf = &oneOf
	} else if getMappingKey(node, "object") != nil {
		object := Object{}
		err := node.DecodeWith(decodeLooze, &object)
		if err != nil {
			return err
		}
		model.Object = &object
	} else {
		return yamlError(node, "model should be one of these: object, enum, oneOf; none of these found")
	}

	*value = model
	return nil
}

func (value Model) MarshalYAML() (interface{}, error) {
	var modelValue interface{}
	if value.IsObject() {
		modelValue = value.Object
	} else if value.IsEnum() {
		modelValue = value.Enum
	} else if value.IsOneOf() {
		modelValue = value.OneOf
	} else {
		return nil, errors.New("Unknown model type")
	}
	modelMap := yamlx.Map()
	if value.Description != nil {
		err := modelMap.Add("description", *value.Description)
		if err != nil {
			return nil, err
		}
	}
	err := modelMap.Merge(modelValue)
	if err != nil {
		return nil, err
	}
	return modelMap.Node, nil
}

func unmarshalModel(keyNode *yaml.Node, valueNode *yaml.Node) (*NamedModel, error) {
	name := Name{}
	err := keyNode.DecodeWith(decodeStrict, &name)
	if err != nil {
		return nil, err
	}
	err = name.Check(PascalCase)
	if err != nil {
		return nil, err
	}
	model := Model{}
	err = valueNode.DecodeWith(decodeStrict, &model)
	if err != nil {
		return nil, err
	}
	return &NamedModel{Name: name, Model: model}, nil
}

func (value *Models) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "models should be YAML mapping")
	}
	count := len(node.Content) / 2
	array := Models{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]
		model, err := unmarshalModel(keyNode, valueNode)
		if err != nil {
			return err
		}
		array = append(array, *model)
	}
	*value = array
	return nil
}

func (value Models) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	for index := 0; index < len(value); index++ {
		model := value[index]
		err := yamlMap.Add(model.Name, model.Model)
		if err != nil {
			return nil, err
		}
	}
	return yamlMap.Node, nil
}

type ModelsMap map[string]*NamedModel

func buildModelsMap(models Models) ModelsMap {
	result := make(map[string]*NamedModel)
	for modIndex := range models {
		name := models[modIndex].Name.Source
		result[name] = &models[modIndex]
	}
	return result
}
