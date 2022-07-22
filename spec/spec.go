package spec

import (
	"errors"

	"github.com/specgen-io/specgen/yamlx/v2"
	"gopkg.in/specgen-io/yaml.v3"
)

type Spec struct {
	Meta
	Versions []Version
}

type VersionSpecification struct {
	Http   Apis   `yaml:"http"`
	Models Models `yaml:"models"`
}

type Version struct {
	Version Name
	VersionSpecification
	ResolvedModels []*NamedModel
}

type Meta struct {
	SpecVersion string  `yaml:"spec"`
	Name        Name    `yaml:"name"`
	Title       *string `yaml:"title,omitempty"`
	Description *string `yaml:"description,omitempty"`
	Version     string  `yaml:"version"`
}

func isVersionNode(node *yaml.Node) bool {
	return VersionFormat.Check(node.Value) == nil
}

func (value *Spec) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return yamlError(node, "http should be YAML mapping")
	}
	versions := []Version{}
	count := len(node.Content) / 2
	for index := 0; index < count; index++ {
		keyNode := node.Content[index*2]
		valueNode := node.Content[index*2+1]

		if isVersionNode(keyNode) {
			version := Name{}
			err := keyNode.DecodeWith(decodeStrict, &version)
			if err != nil {
				return err
			}
			err = version.Check(VersionFormat)
			if err != nil {
				return err
			}

			theSpec := VersionSpecification{}
			valueNode.DecodeWith(decodeStrict, &theSpec)
			versions = append(versions, Version{version, theSpec, nil})
		}
	}
	theSpec := VersionSpecification{}
	err := node.DecodeWith(decodeLooze, &theSpec)
	if err != nil {
		return err
	}
	versions = append(versions, Version{Name{}, theSpec, nil})

	meta := Meta{}
	node.DecodeWith(decodeStrict, &meta)

	*value = Spec{meta, versions}
	return nil
}

func (value VersionSpecification) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	if len(value.Http.Apis) > 0 {
		yamlMap.Add("http", value.Http)
	}
	if len(value.Models) > 0 {
		yamlMap.Add("models", value.Models)
	}
	return yamlMap.Node, nil
}

func (value Spec) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Merge(value.Meta)
	for index := 0; index < len(value.Versions); index++ {
		version := value.Versions[index]

		if version.Version.Source != "" {
			yamlMap.Add(version.Version.Source, version.VersionSpecification)
		} else {
			yamlMap.Merge(version.VersionSpecification)
		}
	}
	return yamlMap.Node, nil
}

func (value Meta) MarshalYAML() (interface{}, error) {
	yamlMap := yamlx.Map()
	yamlMap.Add("spec", yamlx.String(value.SpecVersion))
	yamlMap.Add("name", value.Name)
	yamlMap.AddOmitNil("title", value.Title)
	yamlMap.AddOmitNil("description", value.Description)
	yamlMap.Add("version", yamlx.String(value.Version))
	return yamlMap.Node, nil
}

func unmarshalSpec(data []byte) (*Spec, *Messages, error) {
	messages := NewMessages()
	var spec Spec
	err := yaml.UnmarshalWith(decodeStrict, data, &spec)
	if err != nil {
		messages.Add(convertYamlError(err, nil))
		return nil, messages, errors.New("failed to read specification")
	}
	return &spec, messages, nil
}
