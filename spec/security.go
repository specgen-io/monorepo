package spec

import (
	"gopkg.in/specgen-io/yaml.v3"
)

type securityRef struct {
	Name Name
	Args []string
}

type SecurityRef securityRef

func (value *SecurityRef) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "security reference should be YAML mapping")
	}
	keyNode := node.Content[0]
	valueNode := node.Content[1]

	name := Name{}
	err := keyNode.DecodeWith(decodeStrict, &name)
	if err != nil {
		return err
	}

	args := []string{}
	err = valueNode.DecodeWith(decodeStrict, &args)
	if err != nil {
		return err
	}

	*value = SecurityRef{name, args}
	return nil
}

func (value SecurityRef) MarshalYAML() (interface{}, error) {
	return nil, nil
}

type Security []SecurityRef
