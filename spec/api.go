package spec

import (
	"github.com/specgen-io/specgen/v2/yamlx"
	"gopkg.in/specgen-io/yaml.v3"
)

type Api struct {
	Name       Name
	Operations Operations
	Apis       *Apis
}

type Apis struct {
	Url     *string
	Apis    []Api
	Version *Version
}

func (apis *Apis) GetUrl() string {
	if apis.Url != nil {
		return *apis.Url
	}
	if apis.Version.Version.Source != "" {
		return "/" + apis.Version.Version.Source
	}
	return ""
}

func (value *Apis) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "apis should be YAML mapping")
	}

	url, err := decodeStringOptional(node, "url")
	if err != nil {
		return err
	}

	count := len(node.Content) / 2
	array := []Api{}
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		if !contains([]string{"url"}, keyNode) {
			valueNode := node.Content[index*2+1]
			name := Name{}
			err := keyNode.DecodeWith(decodeStrict, &name)
			if err != nil {
				return err
			}
			err = name.Check(SnakeCase)
			if err != nil {
				return err
			}
			operations := Operations{}
			err = valueNode.DecodeWith(decodeStrict, &operations)
			if err != nil {
				return err
			}
			array = append(array, Api{Name: name, Operations: operations})
		}
	}

	*value = Apis{Url: url, Apis: array}
	return nil
}

func (value Apis) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.AddOmitNil("url", value.Url)
	for index := 0; index < len(value.Apis); index++ {
		api := value.Apis[index]
		yamlMap.Add(api.Name, api.Operations)
	}
	return yamlMap.Node, nil
}
