package old

import (
	"errors"
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type Model struct {
	Object *Object
	Enum   *Enum
	OneOf  *OneOf
}

type NamedModel struct {
	Name Name
	Model
	Version *Version
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

func (model *Model) Description() *string {
	if model.IsObject() {
		return model.Object.Description
	}
	if model.IsOneOf() {
		return model.OneOf.Description
	}
	if model.IsEnum() {
		return model.Enum.Description
	}
	return nil
}

func (value *Model) UnmarshalYAML(node *yaml.Node) error {
	model := Model{}

	if getMappingKey(node, "enum") != nil {
		enum := Enum{}
		err := node.DecodeWith(decodeStrict, &enum)
		if err != nil {
			return err
		}
		model.Enum = &enum
	} else if getMappingKey(node, "oneOf") != nil {
		oneOf := OneOf{}
		err := node.DecodeWith(decodeStrict, &oneOf)
		if err != nil {
			return err
		}
		model.OneOf = &oneOf
	} else {
		object := Object{}
		err := node.DecodeWith(decodeStrict, &object)
		if err != nil {
			return err
		}
		model.Object = &object
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
	if value.Description() != nil {
		err := modelMap.Add("description", *value.Description())
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
	if model.IsEnum() && model.Enum.Description == nil {
		model.Enum.Description = getDescription(keyNode)
	}
	if model.IsObject() && model.Object.Description == nil {
		model.Object.Description = getDescription(keyNode)
	}
	if model.IsOneOf() && model.OneOf.Description == nil {
		model.OneOf.Description = getDescription(keyNode)
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
